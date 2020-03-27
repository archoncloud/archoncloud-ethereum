package client_utils

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

// This function is available to client that downloads an archon file,
// so that they may validate the ContainerSignature associated with
// the download
func GetPublickeyFromContract(username string, timeout time.Duration) (pubkey [64]byte, err error) {
	var ret [64]byte
	if len(username) > 32 {
		return ret, fmt.Errorf("error GetPublickeyFromContract, username must be <= 32 chars")
	}
	var bUsername [32]byte
	copy(bUsername[:], []byte(username)[0:32])

	timeoutMessage := make(chan bool, 1)
	go func(s time.Duration) {
		time.Sleep(s * time.Second)
		timeoutMessage <- true
	}(timeout)

	functionCompleteMessage := make(chan bool, 1)
	retMessage := make(chan [64]byte, 1)
	retErrorMessage := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		// setup channels
		retXMessage := make(chan [32]byte, 1)
		retYMessage := make(chan [32]byte, 1)
		retXMessage_err := make(chan error, 1)
		retYMessage_err := make(chan error, 1)
		go func(username [32]byte, wg *sync.WaitGroup) {
			defer wg.Done()
			resX, resX_err := getPublickeyFromContract(username, "x")
			if resX_err != nil {
				var empty [32]byte
				retXMessage <- empty
				retXMessage_err <- resX_err
			}
			retXMessage <- resX
			retXMessage_err <- nil
		}(bUsername, &wg)

		go func(username [32]byte, wg *sync.WaitGroup) {
			defer wg.Done()
			resY, resY_err := getPublickeyFromContract(bUsername, "y")
			if resY_err != nil {
				var empty [32]byte
				retYMessage <- empty
				retYMessage_err <- resY_err
			}
			retYMessage <- resY
			retYMessage_err <- nil
		}(bUsername, &wg)

		wg.Wait()
		// collect data
		var ret [64]byte
		retX_err := <-retXMessage_err
		retY_err := <-retYMessage_err
		if retX_err != nil || retY_err != nil {
			functionCompleteMessage <- true
			retMessage <- ret
			retErrorMessage <- fmt.Errorf("error GetPublickeyFromContract, rpc call failure")
		}
		resX := <-retXMessage
		resY := <-retYMessage
		copy(ret[0:32], resX[:])
		copy(ret[33:64], resY[:])
		functionCompleteMessage <- true
		retMessage <- ret
		retErrorMessage <- nil
	}()

	select {
	case <-timeoutMessage:
		var empty [64]byte
		return empty, fmt.Errorf("GetPublickeyFromContract has timed-out")
	case <-functionCompleteMessage:
		return <-retMessage, nil
	}
	return ret, nil
}

func getPublickeyFromContract(username [32]byte, xOrY string) (pubkeyXOrY [32]byte, err error) {
	// format storage query
	var keyAndSlot [64]byte
	for i := 0; i < 32; i++ {
		keyAndSlot[i] = username[i]
	}
	keyAndSlot[63] = byte(4) // username2UserProfile
	storagePosition := ethcrypto.Keccak256(keyAndSlot[:])
	var idx byte = 0
	if xOrY == "x" {
		idx = 1
	} else if xOrY == "y" {
		idx = 2
	} else {
		var empty [32]byte
		return empty, fmt.Errorf("error getPublicKeyFromContract, pick x or y only")
	}
	storagePosition[31] += idx // gives usersPublicKey(X or Y)
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
		return empty, fmt.Errorf("error getPublicKeyFromContract, username not registered")
	} else {
		var ret []byte
		for i := 2; i < len(response.Result); i += 2 {
			r, _ := hexutil.Decode("0x" + response.Result[i:i+2])
			ret = append(ret, []byte(r)...)
		}
		var bRet [32]byte
		for i := 0; i < 32; i++ {
			bRet[i] = ret[i]
		}
		return bRet, nil
	}
	var empty [32]byte
	return empty, fmt.Errorf("error getPublicKeyFromContract, username not registered")
}
