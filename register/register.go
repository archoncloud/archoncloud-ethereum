package register

import (
	"fmt"
	"math/big"
	"strings"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
	"github.com/archoncloud/archoncloud-ethereum/encodings"
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
	"github.com/archoncloud/archoncloud-ethereum/wallet"

	"github.com/archoncloud/archoncloud-go/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/pariz/gountries"
)

var MaxSLALevel int = 8 // this can be anything
var g_ethRpcUrl string = archonAbi.RpcUrl()
var g_chainID int64 = archonAbi.ChainID() //1 //18
// 3 for ropsten, 18 for ganache. don't need to set for gorli

type SPParams struct {
	Wallet         wallet.EthereumKeyset
	SLALevel       int
	PledgedStorage uint64
	Bandwidth      uint64

	CountryCode gountries.Codes // must contain A2 field
	MinAskPrice uint64          // Wei per MByte

	Stake         uint64
	HardwareProof [32]byte

	NodeID string
}

func encodeCountryCode(alpha2cc gountries.Codes) ([2]byte, error) {
	// encodes using 10 bits
	var ret [2]byte
	if alpha2cc.Alpha2 == "" {
		return ret, fmt.Errorf(
			"registersp.encodeCountryCode error: CountryCode must be ISO Alpha2")
	}
	ALPHA := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes := []rune(alpha2cc.Alpha2)
	a := strings.Index(ALPHA, string(runes[0]))
	b := strings.Index(ALPHA, string(runes[1]))
	if a < 0 || b < 0 {
		return ret, fmt.Errorf(
			"registersp.encodeCountryCode error: CountryCode must be ISO Alpha2")
	}
	ret[0] = byte(a)
	ret[0] |= byte(b) << 5
	ret[1] |= byte(b) >> 3
	return ret, nil
}

func decodeCountryCode(tenBits [2]byte) gountries.Codes {
	var ret gountries.Codes
	idx0 := int(tenBits[0] & 0x1F)
	idx1 := int((tenBits[0]&0xE0)>>5) + int((tenBits[1]&0x03)<<3)
	ALPHA := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes := []rune(ALPHA)
	a0 := string(runes[idx0])
	a1 := string(runes[idx1])
	cc := a0 + a1
	ret.Alpha2 = string(cc)
	return ret
}

func encodeParams(params SPParams) (res [32]byte, err error) {
	var ret [32]byte
	// SLALevel int (1..3 or so) // reserve 3 bits for up to 8 levels
	if params.SLALevel > 0 && params.SLALevel <= MaxSLALevel {
		ret[0] = byte(params.SLALevel)
	} else {
		var empty [32]byte
		return empty, fmt.Errorf(
			"registersp.encodeParams error: Invalid SLALevel! Range is [1, 8]")
	}
	// PledgedStorage uint64
	ret[1] = byte((params.PledgedStorage & 0xFF00000000000000) >> 56)
	ret[2] = byte((params.PledgedStorage & 0x00FF000000000000) >> 48)
	ret[3] = byte((params.PledgedStorage & 0x0000FF0000000000) >> 40)
	ret[4] = byte((params.PledgedStorage & 0x000000FF00000000) >> 32)
	ret[5] = byte((params.PledgedStorage & 0x00000000FF000000) >> 24)
	ret[6] = byte((params.PledgedStorage & 0x0000000000FF0000) >> 16)
	ret[7] = byte((params.PledgedStorage & 0x000000000000FF00) >> 8)
	ret[8] = byte((params.PledgedStorage & 0x00000000000000FF) >> 0)

	// Bandwidth uint64
	ret[9] = byte((params.Bandwidth & 0xFF00000000000000) >> 56)
	ret[10] = byte((params.Bandwidth & 0x00FF000000000000) >> 48)
	ret[11] = byte((params.Bandwidth & 0x0000FF0000000000) >> 40)
	ret[12] = byte((params.Bandwidth & 0x000000FF00000000) >> 32)
	ret[13] = byte((params.Bandwidth & 0x00000000FF000000) >> 24)
	ret[14] = byte((params.Bandwidth & 0x0000000000FF0000) >> 16)
	ret[15] = byte((params.Bandwidth & 0x000000000000FF00) >> 8)
	ret[16] = byte((params.Bandwidth & 0x00000000000000FF) >> 0)

	// MinAskPrice uint64
	ret[17] = byte((params.MinAskPrice & 0xFF00000000000000) >> 56)
	ret[18] = byte((params.MinAskPrice & 0x00FF000000000000) >> 48)
	ret[19] = byte((params.MinAskPrice & 0x0000FF0000000000) >> 40)
	ret[20] = byte((params.MinAskPrice & 0x000000FF00000000) >> 32)
	ret[21] = byte((params.MinAskPrice & 0x00000000FF000000) >> 24)
	ret[22] = byte((params.MinAskPrice & 0x0000000000FF0000) >> 16)
	ret[23] = byte((params.MinAskPrice & 0x000000000000FF00) >> 8)
	ret[24] = byte((params.MinAskPrice & 0x00000000000000FF) >> 0)

	encCC, err := encodeCountryCode(params.CountryCode)
	if err != nil {
		var empty [32]byte
		return empty, err
	}
	ret[25] |= encCC[0]
	ret[26] |= encCC[1]
	// NOTE THERE ARE SIX BITS AVAILABLE IN IDX 26 SINCE encCC is only 10 bits
	// NOTE IDX's 27 to 31 are available

	return ret, nil
}

func DecodeParams(params [32]byte) (res *SPParams) {
	ret := new(SPParams)
	// SLALEVEL
	ret.SLALevel = int(params[0])
	// AVAILABLE STORAGE
	ret.PledgedStorage = 0
	for i := 0; i < 8; i++ {
		ret.PledgedStorage += uint64(params[i+1]) << uint((7-i)*8)
	}
	// BANDWIDTH
	ret.Bandwidth = 0
	for i := 0; i < 8; i++ {
		ret.Bandwidth += uint64(params[i+9]) << uint((7-i)*8)
	}
	// MINASKPRICE
	ret.MinAskPrice = 0
	for i := 0; i < 8; i++ {
		ret.MinAskPrice += uint64(params[i+17]) << uint((7-i)*8)
	}
	var bCC [2]byte
	copy(bCC[:], params[25:27])
	decCC := decodeCountryCode(bCC)
	ret.CountryCode = decCC

	return ret
}

func CheckIfAddressIsRegistered(address common.BCAddress) (bool, error) {
	var bAddress []byte
	sAddress := string(address)
	for i := 0; i < len(address); i += 2 {
		r, _ := hexutil.Decode("0x" + sAddress[i:i+2])
		bAddress = append(bAddress, []byte(r)...)
	}
	var b20Address [20]byte
	for i := 0; i < 20; i++ {
		b20Address[i] = bAddress[i]
	}
	return CheckIfAddressIsRegistered_byteAddress(b20Address)
}

func CheckIfAddressIsRegistered_byteAddress(address [20]byte) (res bool, err error) {
	var keyAndSlot [64]byte
	for i := 0; i < 20; i++ {
		keyAndSlot[i+12] = address[i]
	}
	keyAndSlot[63] = byte(3) // spAddress2SPProfile

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return false, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func CheckIfInGoodStanding(address common.BCAddress) (bool, error) {
	var bAddress []byte
	sAddress := string(address)
	for i := 0; i < len(address); i += 2 {
		r, _ := hexutil.Decode("0x" + sAddress[i:i+2])
		bAddress = append(bAddress, []byte(r)...)
	}
	var b20Address [20]byte
	for i := 0; i < 20; i++ {
		b20Address[i] = bAddress[i]
	}
	return CheckIfInGoodStanding_byteAddress(b20Address)
}

func CheckIfInGoodStanding_byteAddress(address [20]byte) (res bool, err error) {
	var keyAndSlot [64]byte
	for i := 0; i < 20; i++ {
		keyAndSlot[i+12] = address[i]
	}
	keyAndSlot[63] = byte(3) // spAddress2SPProfile

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(6) // inGoodStanding
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return false, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		return false, nil
	} else if resInt.Text(10) == "1" {
		return true, nil
	}
	return false, fmt.Errorf("error CheckIfInGoodStanding: this case not covered")
}

func checkIfPmtIsEnoughForRegTx(params SPParams) bool {
	// THESE VALUES ARE SAME AS THOSE IN SC
	// DO PAY MIND IF OWNER CHANGES THEIR VALUES
	registerCost := uint64(100000000000) // goal 1e15
	bRegCost := new(big.Int)
	bRegCost.SetUint64(registerCost)
	registerCostScalar := uint64(10000) // goal 1e15
	bRegCostScalar := new(big.Int)
	bRegCostScalar.SetUint64(registerCostScalar)

	zero := uint64(0)
	bZero := new(big.Int)
	bZero.SetUint64(zero)

	slaLevelScalar := bZero       // for now
	pledgedStorageScalar := bZero // for now
	bandwidthScalar := bZero      // for now
	minAskPriceScalar := bZero    // for now

	totalCost := new(big.Int)
	totalCost.Add(totalCost, bRegCost.Mul(bRegCost, bRegCostScalar))
	// slaLevel
	bSLALevel := new(big.Int)
	bSLALevel.SetUint64(uint64(params.SLALevel))
	totalCost.Add(totalCost, bSLALevel.Mul(bSLALevel, slaLevelScalar))

	// avaStorage
	bPledgedStorage := new(big.Int)
	bPledgedStorage.SetUint64(params.PledgedStorage)
	totalCost.Add(totalCost, bPledgedStorage.Mul(bPledgedStorage, pledgedStorageScalar))

	// bandWidth
	bBandwidth := new(big.Int)
	bBandwidth.SetUint64(params.Bandwidth)
	totalCost.Add(totalCost, bBandwidth.Mul(bBandwidth, bandwidthScalar))

	// minAskPrice
	bMinAskPrice := new(big.Int)
	bMinAskPrice.SetUint64(params.MinAskPrice)
	totalCost.Add(totalCost, bMinAskPrice.Mul(bMinAskPrice, minAskPriceScalar))

	// CHECK
	difference := new(big.Int)
	difference.SetUint64(params.Stake)

	difference.Sub(difference, totalCost)
	if difference.Sign() < 0 {
		return false
	}
	return true
}

func RegisterSP(params SPParams) (txid string, err error) {
	encodedParams, c_err := encodeParams(params)
	if c_err != nil {
		return "", c_err
	}

	// construct tx
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
	amount := new(big.Int)
	amount.SetUint64(params.Stake)
	isEnough := checkIfPmtIsEnoughForRegTx(params)
	if !isEnough {
		return "", fmt.Errorf("error RegisterSP: stake pmt is not enough to cover registration fee. See SC for total cost of registration")
	}
	height, height_err := rpc_utils.GetBlockHeight()
	if height_err != nil {
		return "", height_err
	}
	gasLimit, gl_err := rpc_utils.GetGasLimit(height)
	if gl_err != nil {
		return "", gl_err
	}
	r := strings.NewReader(archonAbi.Abi)
	scAbi, err_scAbi := abi.JSON(r) // reader io.Reader
	if err_scAbi != nil {
		return "", err_scAbi
	}

	bNodeID := []byte(params.NodeID)
	var b32NodeID [32]byte
	copy(b32NodeID[:], bNodeID[2:])

	args := encodings.ArchonSCArgs{Params: encodedParams,
		NodeID:        b32NodeID,
		HardwareProof: params.HardwareProof}
	var methodName string = "registerSP"
	dataFormatted, df_err := encodings.FormatData(scAbi, methodName, args)
	if df_err != nil {
		return "", df_err
	}
	contractAddressString := strings.Replace(archonAbi.ContractAddress(), "0x", "", 1)
	var contractAddress []byte
	for i := 0; i < len(contractAddressString); i += 2 {
		r, _ := hexutil.Decode("0x" + contractAddressString[i:i+2])
		contractAddress = append(contractAddress, []byte(r)...)
	}
	var bContractAddress [20]byte
	for i := 0; i < 20; i++ {
		bContractAddress[i] = contractAddress[i]
	}
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Wallet.Address, bContractAddress, amount, dataFormatted)
	if gp_err != nil {
		return "", gp_err
	}
	accountHasEnoughEthers, balance, totalCost, err := rpc_utils.CheckTxCostAgainstBalance(params.Stake, gasLimit, params.Wallet.Address)
	if err != nil {
		return "", err
	}
	if !accountHasEnoughEthers {
		return "", fmt.Errorf("error RegisterSP: totalCost of tx is ", totalCost, " but account balance is ", balance)
	}
	tx := types.NewTransaction(nonce,
		bContractAddress,
		amount,
		gasLimit,
		gasPrice,
		dataFormatted)
	// sign tx
	signedTx, err_signedTx := params.Wallet.SignTx(tx)
	if err_signedTx != nil {
		return "", err_signedTx
	}

	// send tx
	txidString, txid_err := rpc_utils.SendRawTx(signedTx)
	if txid_err != nil {
		return "", txid_err
	}
	return txidString, nil
}

func UnregisterSP(params SPParams) (ret string, err error) {
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
	amount := new(big.Int)
	amount.SetUint64(uint64(0))
	height, h_err := rpc_utils.GetBlockHeight()
	if h_err != nil {
		return "", h_err
	}
	gasLimit, gl_err := rpc_utils.GetGasLimit(height)
	if gl_err != nil {
		return "", gl_err
	}
	r := strings.NewReader(archonAbi.Abi)
	scAbi, err_scAbi := abi.JSON(r) // reader io.Reader
	if err_scAbi != nil {
		return "", err_scAbi
	}

	args := encodings.ArchonSCArgs{}
	var methodName string = "unregisterSP"
	dataFormatted, df_err := encodings.FormatData(scAbi, methodName, args)
	if df_err != nil {
		return "", df_err
	}

	contractAddressString := strings.Replace(archonAbi.ContractAddress(), "0x", "", 1)
	var contractAddress []byte
	for i := 0; i < len(contractAddressString); i += 2 {
		r, _ := hexutil.Decode("0x" + contractAddressString[i:i+2])
		contractAddress = append(contractAddress, []byte(r)...)
	}
	var bContractAddress [20]byte
	for i := 0; i < 20; i++ {
		bContractAddress[i] = contractAddress[i]
	}
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Wallet.Address, bContractAddress, amount, dataFormatted)
	if gp_err != nil {
		return "", gp_err
	}
	accountHasEnoughEthers, balance, totalCost, err := rpc_utils.CheckTxCostAgainstBalance(uint64(0), gasLimit, params.Wallet.Address)
	if err != nil {
		return "", err
	}
	if !accountHasEnoughEthers {
		return "", fmt.Errorf("error RegisterSP: totalCost of tx is ", totalCost, " but account balance is ", balance)
	}
	tx := types.NewTransaction(nonce,
		bContractAddress,
		amount,
		gasLimit,
		gasPrice,
		dataFormatted)
	// sign tx
	signedTx, err_signedTx := params.Wallet.SignTx(tx)
	if err_signedTx != nil {
		return "", err_signedTx
	}

	// send tx
	unregisterSPTxId, u_err := rpc_utils.SendRawTx(signedTx)
	if u_err != nil {
		return "", u_err
	}
	return unregisterSPTxId, nil
}
