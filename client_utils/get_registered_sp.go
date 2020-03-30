package client_utils

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/archoncloud/archoncloud-ethereum/encodings"
	"github.com/archoncloud/archoncloud-ethereum/register"
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"

	"github.com/archoncloud/archoncloud-go/common"
	. "github.com/archoncloud/blockchainAPI/registered_sp"
)

// An sp participating in archon network makes available to archon clients
// its cache of registeredSp profiles. The sp fills this cache by taking
// census of the nodeIDs of its "neighbors" in the overlay, then retrieving
// the sp profile associated with each of thes nodeIDs
func GetRegisteredSP(ethAddress [20]byte) (sp *RegisteredSp, err error) {
	rpcResult, err_result := getRegisteredSpRpcCalls(ethAddress)
	if err_result != nil {
		empty := new(RegisteredSp)
		return empty, err_result
	}
	ret, err_ret := rpcResultToRegisteredSp(*rpcResult)
	if err_ret != nil {
		empty := new(RegisteredSp)
		return empty, err_result
	}
	ret.Address = EthAddressToBCAddress(ethAddress)
	return ret, nil
}

type getRegisteredSpRpcResult struct {
	Params           [32]byte
	NodeID           [32]byte
	Stake            uint64
	RemainingStorage uint64
}

func BCAddressToEthAddress(address common.BCAddress) [20]byte {
	var bAddress []byte
	sAddress := string(address)
	for i := 0; i < len(address); i += 2 {
		r, _ := hexutil.Decode("0x" + sAddress[i:i+2])
		bAddress = append(bAddress, []byte(r)...)
	}
	var b20Address [20]byte
	copy(b20Address[0:20], bAddress[0:20])
	return b20Address
}

func EthAddressToBCAddress(address [20]byte) common.BCAddress {
	var bAddress []byte
	bAddress = make([]byte, 20)
	copy(bAddress[0:len(address)], address[0:len(address)])
	hexAddress := hexutil.Encode(bAddress)
	ret := common.BCAddress(hexAddress)
	return ret
}

func rpcCallParams(ethAddress [20]byte) (res [32]byte, err error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], ethAddress[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile
	// params is first index so dont need to increment

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
		return empty, fmt.Errorf("error rpcCallParams, sp not registered")
	} else {
		var ret []byte
		for i := 2; i < len(response.Result); i += 2 {
			r, _ := hexutil.Decode("0x" + response.Result[i:i+2])
			ret = append(ret, []byte(r)...)
		}
		for len(ret) < 32 {
			ret = append(ret, byte(0))
		}
		var bRet [32]byte
		copy(bRet[0:32], ret[0:32])
		return bRet, nil
	}
	var empty [32]byte
	return empty, fmt.Errorf("error rpcCallParams, sp not registered")
}

func rpcCallNodeID(ethAddress [20]byte) (res [32]byte, err error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], ethAddress[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile

	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(1) // nodeID
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
		return empty, fmt.Errorf("error rpcCallNodeID, sp not registered")
	} else {
		var ret []byte
		for i := 2; i < len(response.Result); i += 2 {
			r, _ := hexutil.Decode("0x" + response.Result[i:i+2])
			ret = append(ret, []byte(r)...)
		}
		for len(ret) < 32 {
			ret = append(ret, byte(0))
		}
		var bRet [32]byte
		copy(bRet[0:32], ret[0:32])
		return bRet, nil
	}
	var empty [32]byte
	return empty, fmt.Errorf("error rpcCallNodeID, sp not registered")
}

func rpcCallStake(ethAddress [20]byte) (res uint64, err error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], ethAddress[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile
	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(2) // stake
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return uint64(0), err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	if resInt.Text(10) == "0" {
		return uint64(0), fmt.Errorf("error rpcCallStake, sp not registered")
	} else {
		stake, stake_err := strconv.ParseUint(strings.Replace(response.Result, "0x", "", 1), 16, 64)
		if stake_err != nil {
			return uint64(0), stake_err
		}
		return stake, nil
	}
	return uint64(0), fmt.Errorf("error rpcCallStake, sp not registered")
}

func rpcCallGetRemainingStorage(ethAddress [20]byte) (res uint64, err error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], ethAddress[0:20])
	keyAndSlot[63] = byte(3) // spAddress2SPProfile
	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	storagePosition[31] += byte(5) // remainingStorage
	hexStoragePosition := hexutil.Encode(storagePosition)
	response, err := rpc_utils.GetStorageAt(hexStoragePosition)
	if err != nil {
		return uint64(0), err
	}
	resInt := new(big.Int)
	resInt.SetString(response.Result[2:], 16)
	remainingStorage, err := strconv.ParseUint(strings.Replace(response.Result, "0x", "", 1), 16, 64)
	if err != nil {
		return uint64(0), err
	}
	return remainingStorage, nil
}

func rpcCallGetLogs(topics []string, fromBlock, toBlock string) (TxLogs, error) {
	logs, err := rpc_utils.GetLogs(topics, fromBlock, toBlock)
	if err != nil {
		var empty TxLogs
		return empty, err
	}
	return TxLogs(logs), nil
}

func getRegisteredSpRpcCalls(ethAddress [20]byte) (res *getRegisteredSpRpcResult,
	err error) {

	ret := new(getRegisteredSpRpcResult)

	inGoodStanding, err := register.CheckIfInGoodStanding_byteAddress(ethAddress)
	if err != nil {
		return ret, err
	}
	if !inGoodStanding {
		return ret, fmt.Errorf("error getRegisteredSpRpcCalls: sp not in good standing")
	}
	var wg sync.WaitGroup
	paramsMessage := make(chan [32]byte, 1)
	paramsMessage_err := make(chan error, 1)
	nodeIDMessage := make(chan [32]byte, 1)
	nodeIDMessage_err := make(chan error, 1)
	stakeMessage := make(chan uint64, 1)
	stakeMessage_err := make(chan error, 1)
	remainingStorageMessage := make(chan uint64, 1)
	remainingStorageMessage_err := make(chan error, 1)
	wg.Add(4)
	// 1.
	go func(ethAddress [20]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		params, err_params := rpcCallParams(ethAddress)
		paramsMessage <- params
		paramsMessage_err <- err_params
	}(ethAddress, &wg)
	// 2.
	go func(ethAddress [20]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		nodeIDEncoded, err_uplUrl := rpcCallNodeID(ethAddress)
		nodeIDMessage <- nodeIDEncoded
		nodeIDMessage_err <- err_uplUrl
	}(ethAddress, &wg)
	// 3.
	go func(ethAddress [20]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		stake, err_stake := rpcCallStake(ethAddress)
		stakeMessage <- stake
		stakeMessage_err <- err_stake
	}(ethAddress, &wg)
	// 4.
	go func(ethAddress [20]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		remainingStorage, err := rpcCallGetRemainingStorage(ethAddress)
		remainingStorageMessage <- remainingStorage
		remainingStorageMessage_err <- err
	}(ethAddress, &wg)

	wg.Wait()
	params_err := <-paramsMessage_err
	nodeID_err := <-nodeIDMessage_err
	stake_err := <-stakeMessage_err
	remainingStorage_err := <-remainingStorageMessage_err
	var errorArray []string
	if params_err != nil {
		errorArray = append(errorArray, params_err.Error())
	}
	if nodeID_err != nil {
		errorArray = append(errorArray, nodeID_err.Error())
	}
	if stake_err != nil {
		errorArray = append(errorArray, stake_err.Error())
	}
	if remainingStorage_err != nil {
		errorArray = append(errorArray, remainingStorage_err.Error())
	}
	if len(errorArray) > 0 {
		empty := new(getRegisteredSpRpcResult)
		return empty, fmt.Errorf(strings.Join(errorArray, "\n"))
	}
	ret.Params = <-paramsMessage
	ret.NodeID = <-nodeIDMessage
	ret.Stake = <-stakeMessage
	ret.RemainingStorage = <-remainingStorageMessage
	return ret, nil
}

func rpcResultToRegisteredSp(rpcResult getRegisteredSpRpcResult) (sp *RegisteredSp,
	err error) {
	var wg sync.WaitGroup
	spParamsMessage := make(chan encodings.SPParams, 1)
	nodeIDMessage := make(chan string, 1)
	wg.Add(2)
	// calls to go routine
	// 1.
	go func(params [32]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		res := encodings.DecodeParams(params) // encodings.SPParams
		spParamsMessage <- *res
	}(rpcResult.Params, &wg)
	// 2.
	go func(bNodeID [32]byte, wg *sync.WaitGroup) {
		defer wg.Done()
		nodeID := string([]byte(bNodeID[:]))
		nodeIDMessage <- nodeID
	}(rpcResult.NodeID, &wg)
	wg.Wait()

	ret := new(RegisteredSp)
	spParams := <-spParamsMessage
	if spParams.MinAskPrice > 0 {
		ret.IsInMarketPlace = true
		ret.MinAskPrice = spParams.MinAskPrice
	} else {
		ret.IsInMarketPlace = false
	}
	ret.SLALevel = spParams.SLALevel
	ret.PledgedStorage = spParams.PledgedStorage
	ret.Bandwidth = spParams.Bandwidth
	ret.CountryCode = spParams.CountryCode
	nodeID := <-nodeIDMessage
	ret.NodeID = nodeID
	ret.RemainingStorage = rpcResult.RemainingStorage

	return ret, nil
}
