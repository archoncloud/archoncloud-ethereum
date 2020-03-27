package rpc_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
)

var g_ethRpcUrl string = archonAbi.RpcUrl()

type Response struct {
	Result string `json:"result"`
}

func GetStorageAt(hexStoragePosition string) (Response, error) {
	var blockParameter string = "latest"
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getStorageAt\",\"params\": [\"" + archonAbi.ContractAddress() + "\", \"" + hexStoragePosition + "\", \"" + blockParameter + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	req, err_req := http.NewRequest("POST", g_ethRpcUrl, bytes.NewBuffer(reqBytes))
	if err_req != nil {
		return *new(Response), fmt.Errorf("error GetStorageAt, network error 1")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err_resp := client.Do(req)
	if resp == nil || err_resp != nil {
		return *new(Response), fmt.Errorf("error GetStorageAt, network error 2")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response Response
	err_json := json.Unmarshal(body, &response)
	if err_json != nil {
		return *new(Response), fmt.Errorf("error GetStorageAt, json parse error")
	}
	return response, nil
}
