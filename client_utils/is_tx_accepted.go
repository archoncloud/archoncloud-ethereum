package client_utils

// Note this is designed for registerSP and proposeUpload tx only. Other txs
// will not necessarily return correct bool.
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
