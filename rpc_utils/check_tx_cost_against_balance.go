package rpc_utils

import (
	"math/big"
)

// This is useful during tx construction
func CheckTxCostAgainstBalance(amount, gasLimit uint64, address [20]byte) (accountHasEnough bool, balance, totalCost big.Int, err error) {
	zero := *new(big.Int)
	zero.SetUint64(uint64(0))
	bBalance, err := GetBalance_byteAddressToBigInt(address)
	if err != nil {
		return false, bBalance, zero, err
	}
	bAmount := new(big.Int)
	bAmount.SetUint64(amount)
	bGasLimit := new(big.Int)
	bGasLimit.SetUint64(gasLimit)

	bTotalCost := new(big.Int)
	bTotalCost = bAmount.Add(bAmount, bGasLimit)

	balanceCopy := new(big.Int)
	balanceCopy.SetString(bBalance.Text(16), 16)

	difference := bBalance.Sub(&bBalance, bTotalCost)
	if difference.Sign() < 0 {
		return false, *balanceCopy, *bTotalCost, nil
	}
	return true, bBalance, *bTotalCost, nil
}
