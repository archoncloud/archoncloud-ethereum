package client_utils

import (
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
)

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
