package client_utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

func GetUsernameFromContract(address [20]byte) (username [32]byte, err error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], address[0:20])
	keyAndSlot[63] = byte(5) // address2Username

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		var empty [32]byte
		return empty, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		var empty [32]byte
		return empty, fmt.Errorf("error GetUsernameFromContract, username not registered")
	} else {
		var ret []byte
		for i := 2; i < len(response.Result); i += 2 {
			r, _ := hexutil.Decode("0x" + response.Result[i:i+2])
			ret = append(ret, []byte(r)...)
		}
		var bRet [32]byte
		copy(bRet[0:32], ret[0:32])
		return bRet, nil
	}
	var empty [32]byte
	return empty, fmt.Errorf("error GetUsernameFromContract, username not registered")
}
