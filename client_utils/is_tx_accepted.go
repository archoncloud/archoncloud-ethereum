package client_utils

func IsTxAcceptedByBlockchain(txid string) (bool, error) {
	logs, err := GetTxLogs(txid)
	if err != nil {
		return false, err
	}
	if len(logs) < 1 {
		return false, err
	}
	return true, nil
}
