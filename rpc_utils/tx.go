package rpc_utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction
	"github.com/ethereum/go-ethereum/rlp"
)

func GetNonceForAddress(address [20]byte) (ret uint64, err error) {
	var bAddress []byte
	bAddress = append(bAddress, address[:]...)
	hexAddress := hexutil.Encode(bAddress)
	blockParameter := "pending"
	var reqString = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\": [\"" + hexAddress + "\",\"" + blockParameter + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return uint64(0), err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return uint64(0), err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Response struct {
		Result string `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return uint64(0), json_err
	}
	res, err_res := strconv.ParseUint(strings.Replace(response.Result, "0x", "", 1),
		16,
		64)
	if err_res != nil {
		return uint64(0), err_res
	}
	return uint64(res), nil
}

func GetBlockHeight() (ret string, err error) {
	var reqBytes = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return "", err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return "", err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Response struct {
		Result string `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return "", json_err
	}

	return response.Result, nil
}

func GetBlockHash(height string) (string, error) {

	var reqBytes = []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["` + height + `", false],"id":1}`)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return "", err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return "", err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Block struct {
		Hash string `json:"hash"`
	}
	type Response struct {
		Result Block `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return "", json_err
	}
	return response.Result.Hash, nil
}

func GetGasLimit(height string) (res uint64, err error) {
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\": [\"" + height + "\",false],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return uint64(0), err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return uint64(0), err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type ResultStruct struct {
		GasLimit string `json:"gasLimit"`
	}
	type Response struct {
		Result ResultStruct `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return uint64(0), json_err
	}

	ret, err_ret := hexutil.DecodeUint64(response.Result.GasLimit)
	if err_ret != nil {
		return uint64(0), err_ret
	}
	return ret, nil
}

func EstimateGas(from [20]byte, to [20]byte, value *big.Int, data []byte) (res *big.Int, err error) {
	var bFrom []byte
	bFrom = append(bFrom, from[:]...)
	hexFrom := hexutil.Encode(bFrom)
	var bTo []byte
	bTo = append(bTo, from[:]...)
	hexTo := hexutil.Encode(bTo)
	hexValue := hexutil.EncodeBig(value)
	hexData := hexutil.Encode(data)
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_estimateGas\",\"params\": [{\"from\": \"" + hexFrom + "\",\"to\": \"" + hexTo + "\",\"value\": \"" + hexValue + "\",\"data\": \"" + hexData + "\"}],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return new(big.Int), err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return new(big.Int), err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Response struct {
		Result string `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return new(big.Int), json_err
	}
	bRes := hexutil.MustDecode(response.Result)
	ret := new(big.Int)
	ret.SetBytes(bRes)

	return ret, nil
}

func SendRawTx(signedTx *types.Transaction) (res string, err error) {
	data, err_data := rlp.EncodeToBytes(signedTx)
	if err_data != nil {
		return "", err_data
	}
	hexData := hexutil.Encode(data)
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendRawTransaction\",\"params\":[\"" + hexData + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return "", err_req
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return "", err_resp
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	type Response struct {
		Result string `json:"result"`
	}
	var response Response
	json_err := json.Unmarshal(body, &response)
	if json_err != nil {
		return "", json_err
	}
	return response.Result, nil
}
