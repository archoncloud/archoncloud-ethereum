package rpc_utils

import (
	"fmt"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
)

func SanitizeRpcUrls(rpcPtr *archonAbi.UrlPtr) error {
	var updatedUrls []string
	for i := 0; i < len(rpcPtr.Urls); i++ {
		rpcPtr.Url = rpcPtr.Urls[i]
		_, err := GetBlockHeight()
		if err == nil {
			updatedUrls = append(updatedUrls, rpcPtr.Urls[i])
		}
	}
	if len(updatedUrls) < 1 {
		return fmt.Errorf("Error SanitizeRpcUrls: all rpc urls invalid.")
	}
	rpcPtr.Urls = updatedUrls
	rpcPtr.Url = updatedUrls[0]
	return nil
}
