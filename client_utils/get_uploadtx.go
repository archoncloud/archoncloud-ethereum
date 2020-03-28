package client_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
	"github.com/archoncloud/archoncloud-ethereum/encodings"

	"golang.org/x/crypto/sha3"
)

type InputDeconstruct struct {
	HashedArchonFilepath [32]byte
	ContainerSignature   [ethcrypto.SignatureLength]byte
	Params               encodings.ProposeUploadParams
	Shardsize            uint64
	ArchonSPs            [][20]byte
}

type UploadTx struct {
	BlockHash          string `json:blockHash`
	BlockNumber        string `json:blockNumber`
	From               string `json:from`
	Gas                string `json:gas`
	GasPrice           string `json:gasPrice`
	Hash               string `json:hash`
	InputDeconstructed InputDeconstruct
	Input              string `json:input`
	Nonce              string `json:nonce`
	R                  string `json:r`
	S                  string `json:s`
	To                 string `json:to`
	TransactionIndex   string `json:transactionIndex`
	V                  string `json:v`
	AmountPaid         string
	PublicKey          [64]byte
	UsernameCompressed [32]byte
}

// Called by sp. When uploader sends shards (and associated txid)
// to an sp, the sp checks the data and validity of the tx associated
// with the txid
func GetUploadTx(txHash [32]byte) (uploadTx UploadTx, err error) {
	var wg sync.WaitGroup
	// make rpc call(s)
	data, err_data := checkUploadRpcCalls(txHash)
	if err_data != nil {
		return *new(UploadTx), err_data
	}
	if len(data.Hash) == 0 {
		return *new(UploadTx), fmt.Errorf("GetUploadTx error: tx doesn't exist")
	}
	encodedUsernameMessage := make(chan [32]byte, 1)
	encodedUsernameMessage_err := make(chan error, 1)
	wg.Add(1)
	go func(data GetTxByHashResult, wg *sync.WaitGroup) {
		defer wg.Done()
		address := strings.Replace(data.From, "0x", "", 1)
		var bAddress []byte
		for i := 0; i < len(address); i += 2 {
			r, _ := hexutil.Decode("0x" + address[i:i+2])
			bAddress = append(bAddress, []byte(r)...)
		}
		var bAddress2 [20]byte
		copy(bAddress2[0:20], bAddress[0:20])
		username, err_username := GetUsernameFromContract(bAddress2)
		if err_username != nil {
			var empty [32]byte
			encodedUsernameMessage <- empty
			encodedUsernameMessage_err <- fmt.Errorf("GetUploadTx error: error getting username from contract")
			return
		}
		encodedUsernameMessage <- username
		encodedUsernameMessage_err <- nil
		return
	}(data, &wg)
	r := strings.NewReader(archonAbi.Abi)
	scAbi, err_scAbi := abi.JSON(r) // reader io.Reader
	if err_scAbi != nil {
		return *new(UploadTx), fmt.Errorf("GetUploadTx error: Abi file read error")
	}
	input := strings.Replace(data.Input, "0x", "", 1)
	var bInput []byte
	for i := 0; i < len(input); i += 2 {
		r, _ := hexutil.Decode("0x" + input[i:i+2])
		bInput = append(bInput, []byte(r)...)
	}
	var scMethodName string = "proposeUpload"
	proposeUploadMethod := scAbi.Methods[scMethodName]
	truncatedInput := bInput[len(proposeUploadMethod.ID()):]
	v, err_unpack := Unpack(scMethodName, truncatedInput)
	if err_unpack != nil {
		return *new(UploadTx), fmt.Errorf("GetUploadTx error: Abi file read error")
	}
	pv, ok := v.(ProposeUploadArgs)
	if !ok {
		return *new(UploadTx), fmt.Errorf("GetUploadTx error: ProposeUpload data malformed")
	}
	decoded := encodings.DecodeProposeUploadParams(pv.Params)
	var containerSignature [ethcrypto.SignatureLength]byte
	copy(containerSignature[0:32], pv.ContainerSignatureR[:])
	copy(containerSignature[32:64], pv.ContainerSignatureS[:])
	containerSignature[64] = decoded.ContainerSignatureV
	bPubKey, err_bPubKey := ECRecoverFromTx(data)
	if err_bPubKey != nil {
		return *new(UploadTx), err_bPubKey
	}
	wg.Wait()
	err_usernameMessage := <-encodedUsernameMessage_err
	if err_usernameMessage != nil {
		return *new(UploadTx), err_usernameMessage
	}

	encodedUsername := <-encodedUsernameMessage
	ret := UploadTx{BlockHash: data.BlockHash,
		BlockNumber: data.BlockNumber,
		From:        data.From,
		Gas:         data.Gas,
		GasPrice:    data.GasPrice,
		Hash:        data.Hash,
		InputDeconstructed: InputDeconstruct{ArchonSPs: pv.ArchonSPs,
			HashedArchonFilepath: pv.HashedArchonFilepath,
			Params:               decoded,
			ContainerSignature:   containerSignature,
			Shardsize:            pv.Shardsize},
		Input:              data.Input,
		Nonce:              data.Nonce,
		R:                  data.R,
		S:                  data.S,
		To:                 data.To,
		TransactionIndex:   data.TransactionIndex,
		V:                  data.V,
		AmountPaid:         data.Value,
		PublicKey:          bPubKey,
		UsernameCompressed: encodedUsername}

	return ret, nil
}

type GetTxByHashResult struct {
	BlockHash        string `json:blockHash`
	BlockNumber      string `json:blockNumber`
	From             string `json:from`
	Gas              string `json:gas`
	GasPrice         string `json:gasPrice`
	Hash             string `json:hash`
	Input            string `json:input`
	Nonce            string `json:nonce`
	R                string `json:r`
	S                string `json:s`
	To               string `json:to`
	TransactionIndex string `json:transactionIndex`
	V                string `json:v`
	Value            string `json:value`
}

func checkUploadRpcCalls(txHash [32]byte) (res GetTxByHashResult, err error) {
	// gettxbyid rpc call
	var bTxHash []byte
	bTxHash = append(bTxHash, txHash[:]...)
	hexTxHash := hexutil.Encode(bTxHash)
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionByHash\",\"params\": [\"" + hexTxHash + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", archonAbi.RpcUrl(), bytes.NewBuffer(reqBytes))
	type Response struct {
		Result GetTxByHashResult `json:"result"`
	}
	var response Response
	if err_req != nil {
		return response.Result, fmt.Errorf("Error checkUploadRpcCalls: http request initialization error")
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return response.Result, fmt.Errorf("Error checkUploadRpcCalls: http request error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err_json := json.Unmarshal(body, &response)
	if err_json != nil {
		return response.Result, fmt.Errorf("Error checkUploadRpcCalls: response parsing error")
	}
	return response.Result, nil
}

type ProposeUploadArgs struct {
	HashedArchonFilepath [32]byte
	ContainerSignatureR  [32]byte
	ContainerSignatureS  [32]byte
	Params               [32]byte
	Shardsize            uint64
	ArchonSPs            [][20]byte
}

func Unpack(methodName string, input []byte) (inter interface{}, err error) {
	if methodName == "proposeUpload" {
		if len(input)%32 != 0 && len(input) < (4*32) {
			type Empty struct{}
			return new(Empty), fmt.Errorf("This method cannot be unpacked")
		}
		var hashedArchonFilepath [32]byte
		var containerSignatureR, containerSignatureS [32]byte
		var params [32]byte
		copy(hashedArchonFilepath[0:32], input[0:32])
		copy(containerSignatureR[0:32], input[32:64])
		copy(containerSignatureS[0:32], input[(2*32):(3*32)])
		copy(params[0:32], input[(3*32):(4*32)])
		var shardsize uint64
		for i := 24; i < 32; i++ { // getting the uint64 from the byte32
			var shift = 8 * (32 - i - 1)
			shardsize += uint64(input[i+(4*32)]) << shift
		}

		var archonSPs [][20]byte
		var idx int = 0
		for i := (7 * 32); i < len(input); i += 32 {
			var archonSP [20]byte
			copy(archonSP[0:20], input[i+12:i+12+20])
			archonSPs = append(archonSPs, archonSP)
			idx++
		}
		ret := ProposeUploadArgs{HashedArchonFilepath: hashedArchonFilepath,
			ContainerSignatureR: containerSignatureR,
			ContainerSignatureS: containerSignatureS,
			Params:              params,
			Shardsize:           shardsize,
			ArchonSPs:           archonSPs}
		return ret, nil
	}
	type Empty struct{}
	return new(Empty), fmt.Errorf("This method cannot be unpacked")
}

func ECRecoverFromTx(data GetTxByHashResult) (retKey [64]byte, err error) {
	dataR := strings.Replace(data.R, "0x", "", 1)
	dataS := strings.Replace(data.S, "0x", "", 1)
	dataV := strings.Replace(data.V, "0x", "", 1)
	if len(dataR)%2 == 1 {
		dataR = "0" + dataR
	}
	if len(dataS)%2 == 1 {
		dataS = "0" + dataS
	}
	if len(dataV)%2 == 1 {
		dataV = "0" + dataV
	}
	var sig []byte // = r + s + v
	var bR []byte
	for i := 0; i < len(dataR); i += 2 {
		r, _ := hexutil.Decode("0x" + dataR[i:i+2])
		bR = append(bR, []byte(r)...)
	}
	var bS []byte
	for i := 0; i < len(dataS); i += 2 {
		r, _ := hexutil.Decode("0x" + dataS[i:i+2])
		bS = append(bS, []byte(r)...)
	}
	var bV []byte
	for i := 0; i < len(dataV); i += 2 {
		r, _ := hexutil.Decode("0x" + dataV[i:i+2])
		bV = append(bV, []byte(r)...)
	}
	sig = append(sig, bR[:]...)
	sig = append(sig, bS[:]...)
	sig = append(sig, bV[:]...)

	nonce, err_res := strconv.ParseUint(strings.Replace(data.Nonce, "0x", "", 1), 16, 64)
	if err_res != nil {
		var empty [64]byte
		return empty, fmt.Errorf("ECRecoverFromTx error: parsing nonce error")
	}
	gasPrice, err_gasPrice := strconv.ParseUint(strings.Replace(data.GasPrice, "0x", "", 1), 16, 64)
	if err_gasPrice != nil {
		var empty [64]byte
		return empty, fmt.Errorf("ECRecoverFromTx error: parsing gasPrice error")
	}
	bigGasPrice := new(big.Int)
	bigGasPrice.SetUint64(gasPrice)
	gas, err_gas := strconv.ParseUint(strings.Replace(data.Gas, "0x", "", 1), 16, 64)
	if err_gas != nil {
		var empty [64]byte
		return empty, fmt.Errorf("ECRecoverFromTx error: parsing gas error")
	}
	bigGas := new(big.Int)
	bigGas.SetUint64(gas)
	var bTo []byte
	for i := 2; i < len(data.To); i += 2 {
		r, _ := hexutil.Decode("0x" + data.To[i:i+2])
		bTo = append(bTo, []byte(r)...)
	}
	var bbTo [20]byte
	copy(bbTo[0:20], bTo[0:20])
	dataValue, err_dataValue := strconv.ParseUint(strings.Replace(data.Value, "0x", "", 1), 16, 64)
	if err_dataValue != nil {
		var empty [64]byte
		return empty, fmt.Errorf("ECRecoverFromTx error: parsing value error")
	}
	bigValue := new(big.Int)
	bigValue.SetUint64(dataValue)

	var chainId byte
	if archonAbi.ChainIs() == "Ganorge" {
		chainId = byte(1) // byte(5)
	} else if archonAbi.ChainIs() == "Gorli" {
		chainId = byte(5)
	}
	hw := sha3.NewLegacyKeccak256()
	input := strings.Replace(data.Input, "0x", "", 1)
	var bInput []byte
	for i := 0; i < len(input); i += 2 {
		r, _ := hexutil.Decode("0x" + input[i:i+2])
		bInput = append(bInput, []byte(r)...)
	}
	rlp.Encode(hw, []interface{}{
		nonce,
		bigGasPrice,
		bigGas,
		bbTo,
		bigValue,
		bInput,
		chainId, uint(0), uint(0)})
	var h common.Hash
	hw.Sum(h[:0])
	var bH []byte
	bH = append(bH, []byte(h[:])...)
	if bV[0] == byte(46) || bV[0] == byte(38) {
		sig[64] = byte(1)
	} else {
		sig[64] = byte(0)
	}
	pubKey, err_pub := ethcrypto.Ecrecover(bH, sig)
	if err_pub != nil {
		var empty [64]byte
		return empty, fmt.Errorf("ECRecoverFromTx error: ethcrypto.Ecrecover error")
	}
	var bPubKey [64]byte
	copy(bPubKey[0:64], pubKey[1:65])
	return bPubKey, nil
}
