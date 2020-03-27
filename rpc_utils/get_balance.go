package rpc_utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func GetBalance(hexAddress string) (Response, error) {
	var blockParameter string = "latest"
	var reqString string = "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\": [\"" + hexAddress + "\", \"" + blockParameter + "\"],\"id\":1}"
	var reqBytes = []byte(reqString)
	return HttpPostWResponse(reqBytes)
}

func GetBalance_byteAddress(address [20]byte) (Response, error) {
	var bEthAddress []byte
	bEthAddress = make([]byte, 20)
	copy(bEthAddress[0:20], address[0:20])
	hexAddress := hexutil.Encode(bEthAddress)

	return GetBalance(hexAddress)
}

func GetBalance_byteAddressToBigInt(address [20]byte) (big.Int, error) {
	response, err := GetBalance_byteAddress(address)
	if err != nil {
		return *new(big.Int), err
	}
	var ret big.Int
	ret.SetString(response.Result[2:], 16)
	return ret, nil
}
