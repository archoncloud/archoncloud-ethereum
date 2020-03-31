package rpc_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TxLog struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type TxLogs []TxLog

type TxReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	From              string `json:"from"`
	To                string `json:"to"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress"`
	TxLogs            TxLogs `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
}

type TxReceiptResponse struct {
	Result TxReceipt `json:"result"`
}

func GetTxReceipt(txid string) (TxReceipt, error) {
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionReceipt\",\"params\": [\"" + txid + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpc.Url, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return *new(TxReceipt), fmt.Errorf("error GetTxReceipt, network error 1")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return *new(TxReceipt), fmt.Errorf("error GetTxReceipt, network error 2")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response TxReceiptResponse
	err_json := json.Unmarshal(body, &response)
	if err_json != nil {
		return *new(TxReceipt), fmt.Errorf("error GetTxReceipt, json parse error 1 ", err_json)
	}
	if response.Result.TransactionHash == "" {
		return *new(TxReceipt), fmt.Errorf("error GetTxReceipt, txreceipt null")
	}
	return response.Result, nil
}
