package rpc_utils

import (
	"math/big"
)

// This is useful during tx construction
func CheckTxCostAgainstBalance(amount, gasLimit uint64, address [20]byte) (accountHasEnough bool, balance, totalCost big.Int, err error) {
	bBalance, err := GetBalance_byteAddressToBigInt(address)
	if err != nil {
		return false, bBalance, *new(big.Int), err
	}
	bAmount := new(big.Int)
	bAmount.SetUint64(amount)
	bGasLimit := new(big.Int)
	bGasLimit.SetUint64(gasLimit)

	bTotalCost := bAmount.Add(bAmount, bGasLimit)
	difference := bBalance.Sub(&bBalance, bTotalCost)
	if difference.Sign() < 0 {
		return false, balance, totalCost, nil
	}
	return true, bBalance, *bTotalCost, nil
}
