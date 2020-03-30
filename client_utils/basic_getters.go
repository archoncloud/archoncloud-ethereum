package client_utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// called by CheckTxCostAgainstBalance
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

// utility function for sp to list its earnings
func GetEarnings(ethAddress [20]byte) (big.Int, error) {
	var keyAndSlot [64]byte
	copy(keyAndSlot[12:32], ethAddress[0:20])
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

type TxLog rpc_utils.TxLog
type TxLogs rpc_utils.TxLogs
type TxReceipt rpc_utils.TxReceipt

func GetTxLogs(txid string) (TxLogs, error) {
	receipt, err := rpc_utils.GetTxReceipt(txid)
	if err != nil {
		var empty TxLogs
		return empty, err
	}
	return TxLogs(receipt.TxLogs), nil
}

// An uploader needs to first have their username registered with the sc
// before they can upload to the archon cloud.
// A use of this is that when a downloader needs to validate their download,
// they retrieve from the contract the publicKey associated with the username
// that namespaces the download.
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

// called by routines in archon-dht since in the smart contract
// nodeID -> address -> spProfile
func GetNodeID2Address(nodeID [32]byte) ([20]byte, error) {
	// format storage query
	var keyAndSlot [64]byte
	copy(keyAndSlot[0:32], nodeID[0:32])
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
		copy(bRet[0:20], ret[0:20])
		return bRet, nil
	}
	var empty [20]byte
	return empty, fmt.Errorf("error GetNodeID2Address, nodeID not registered")
}
