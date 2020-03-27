package client_utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

// called by routines in archon-dht since in the smart contract
// nodeID -> address -> spProfile
func GetNodeID2Address(nodeID [32]byte) ([20]byte, error) {
	// format storage query
	var keyAndSlot [64]byte
	for i := 0; i < 32; i++ {
		keyAndSlot[i] = nodeID[i]
	}
	keyAndSlot[63] = byte(6) // nodeID2Address

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		var empty [20]byte
		return empty, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		var empty [20]byte
		return empty, fmt.Errorf("error GetNodeID2Address, nodeID not registered")
	} else {
		var ret []byte
		var startPos int
		// this is to guard against different rpcs returning
		// different paddings for address
		if len(response.Result) == 42 {
			startPos = 0
		} else {
			startPos = 26
		}
		for i := startPos; i < len(response.Result); i += 2 {
			r, _ := hexutil.Decode("0x" + response.Result[i:i+2])
			ret = append(ret, []byte(r)...)
		}
		var bRet [20]byte
		for i := 0; i < 20; i++ {
			bRet[i] = ret[i]
		}
		return bRet, nil
	}
	var empty [20]byte
	return empty, fmt.Errorf("error GetNodeID2Address, nodeID not registered")
}
