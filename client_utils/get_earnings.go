package client_utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

func GetEarnings(ethAddress [20]byte) (big.Int, error) {
	var keyAndSlot [64]byte
	for i := 0; i < 20; i++ {
		keyAndSlot[i+12] = ethAddress[i]
	}
	keyAndSlot[63] = byte(3) // spAddress2SPProfile
	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(3) // 3 earnings
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return *new(big.Int), err
	}
	var ret big.Int
	ret.SetString(response.Result[2:], 16)
	return ret, nil
}
