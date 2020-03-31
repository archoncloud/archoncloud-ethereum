package rpc_utils

import (
	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
)

var g_ethRpc = archonAbi.Rpc()

type Response struct {
	Result string `json:"result"`
}

func GetStorageAt(hexStoragePosition string) (Response, error) {
	var blockParameter string = "latest"
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getStorageAt\",\"params\": [\"" + archonAbi.ContractAddress() + "\", \"" + hexStoragePosition + "\", \"" + blockParameter + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	return HttpPostWResponse(reqBytes)
}
