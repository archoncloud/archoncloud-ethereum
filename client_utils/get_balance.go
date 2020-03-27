package client_utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

func GetBalance(ethAddress [20]byte) (big.Int, error) {
	var bEthAddress []byte
	bEthAddress = make([]byte, 20)
	copy(bEthAddress[0:20], ethAddress[0:20])
	hexAddress := hexutil.Encode(bEthAddress)
	response, err := rpc_utils.GetBalance(hexAddress)
	if err != nil {
		return *new(big.Int), err
	}
	var ret big.Int
	ret.SetString(response.Result[2:], 16)
	return ret, nil
}
