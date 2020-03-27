package rpc_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	archonAbi "github.com/itsmeknt/archoncloud-go/blockchainAPI/ethereum/abi"
)

type LogsResponse struct {
	Result TxLogs `json:"result"`
}

func GetLogs(topics []string, fromBlock, toBlock string) (TxLogs, error) {
	hexContractAddress := archonAbi.ContractAddress
	var sTopics string = "["
	for i := 0; i < len(topics); i++ {
		sTopics += "\"" + topics[i] + "\""
		if i < (len(topics) - 1) {
			sTopics += ","
		}
	}
	sTopics += "]"
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getLogs\",\"params\":[{\"address\": \"" + hexContractAddress() + "\", \"fromBlock\": \"" + fromBlock + "\", \"toBlock\": \"" + toBlock + "\", \"topics\":" + sTopics + "}],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return *new(TxLogs), fmt.Errorf("error GetLogs, network error 1")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return *new(TxLogs), fmt.Errorf("error GetLogs, network error 2")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response LogsResponse
	err_json := json.Unmarshal(body, &response)
	if err_json != nil {
		return *new(TxLogs), fmt.Errorf("error GetLogs, json parse error 1 ", err_json)
	}
	if len(response.Result) < 1 {
		return *new(TxLogs), fmt.Errorf("error GetLogs, txreceipt null")
	}
	return response.Result, nil
}
