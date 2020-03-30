package register

import (
	"fmt"
	"math/big"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func CheckIfAddressIsRegistered(address string) (bool, error) {
	var bAddress []byte
	sAddress := string(address)
	for i := 0; i < len(address); i += 2 {
		r, _ := hexutil.Decode("0x" + sAddress[i:i+2])
		bAddress = append(bAddress, []byte(r)...)
	}
	var b20Address [20]byte
	copy(b20Address[0:20], bAddress[0:20])
	return CheckIfAddressIsRegistered_byteAddress(b20Address)
}

func CheckIfAddressIsRegistered_byteAddress(address [20]byte) (res bool, err error) {
	var keyAndSlot [64]byte
	copy(keyAndslot[12:32], address[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return false, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func CheckIfInGoodStanding(address string) (bool, error) {
	var bAddress []byte
	sAddress := string(address)
	for i := 0; i < len(address); i += 2 {
		r, _ := hexutil.Decode("0x" + sAddress[i:i+2])
		bAddress = append(bAddress, []byte(r)...)
	}
	var b20Address [20]byte
	copy(b20Address[0:20], bAddress[0:20])
	return CheckIfInGoodStanding_byteAddress(b20Address)
}

func CheckIfInGoodStanding_byteAddress(address [20]byte) (res bool, err error) {
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], address[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(6) // inGoodStanding
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return false, err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		return false, nil
	} else if resInt.Text(10) == "1" {
		return true, nil
	}
	return false, fmt.Errorf("error CheckIfInGoodStanding: this case not covered")
}

func checkIfPmtIsEnoughForRegTx(params SPParams) bool {
	// THESE VALUES ARE SAME AS THOSE IN SC
	// DO PAY MIND IF OWNER CHANGES THEIR VALUES
	registerCost := uint64(100000000000) // goal 1e15
	bRegCost := new(big.Int)
	bRegCost.SetUint64(registerCost)
	registerCostScalar := uint64(10000) // goal 1e15
	bRegCostScalar := new(big.Int)
	bRegCostScalar.SetUint64(registerCostScalar)

	zero := uint64(0)
	bZero := new(big.Int)
	bZero.SetUint64(zero)

	slaLevelScalar := bZero       // for now
	pledgedStorageScalar := bZero // for now
	bandwidthScalar := bZero      // for now
	minAskPriceScalar := bZero    // for now

	totalCost := new(big.Int)
	totalCost.Add(totalCost, bRegCost.Mul(bRegCost, bRegCostScalar))
	// slaLevel
	bSLALevel := new(big.Int)
	bSLALevel.SetUint64(uint64(params.SLALevel))
	totalCost.Add(totalCost, bSLALevel.Mul(bSLALevel, slaLevelScalar))

	// avaStorage
	bPledgedStorage := new(big.Int)
	bPledgedStorage.SetUint64(params.PledgedStorage)
	totalCost.Add(totalCost, bPledgedStorage.Mul(bPledgedStorage, pledgedStorageScalar))

	// bandWidth
	bBandwidth := new(big.Int)
	bBandwidth.SetUint64(params.Bandwidth)
	totalCost.Add(totalCost, bBandwidth.Mul(bBandwidth, bandwidthScalar))

	// minAskPrice
	bMinAskPrice := new(big.Int)
	bMinAskPrice.SetUint64(params.MinAskPrice)
	totalCost.Add(totalCost, bMinAskPrice.Mul(bMinAskPrice, minAskPriceScalar))

	// CHECK
	difference := new(big.Int)
	difference.SetUint64(params.Stake)

	difference.Sub(difference, totalCost)
	if difference.Sign() < 0 {
		return false
	}
	return true
}
