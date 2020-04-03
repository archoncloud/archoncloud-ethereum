package register

import (
	"fmt"
	"math/big"
	"strings"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
	"github.com/archoncloud/archoncloud-ethereum/encodings"
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
	"github.com/archoncloud/archoncloud-ethereum/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction

	"github.com/pariz/gountries"
)

type SPParams struct {
	Wallet         wallet.EthereumKeyset
	SLALevel       int
	PledgedStorage uint64
	Bandwidth      uint64

	CountryCode gountries.Codes // must contain A2 field
	MinAskPrice uint64          // Wei per MByte

	Stake         uint64
	HardwareProof [32]byte

	NodeID string
}

func RegisterSP(params SPParams) (txid string, err error) {
	encodedParams, c_err := encodings.EncodeParams(*params.ToEncodingParams())
	if c_err != nil {
		return "", c_err
	}

	// construct tx
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
	amount := new(big.Int)
	amount.SetUint64(params.Stake)
	isEnough := checkIfPmtIsEnoughForRegTx(params)
	if !isEnough {
		return "", fmt.Errorf("error RegisterSP: stake pmt is not enough to cover registration fee. See SC for total cost of registration")
	}
	height, height_err := rpc_utils.GetBlockHeight()
	if height_err != nil {
		return "", height_err
	}
	gasLimit, gl_err := rpc_utils.GetGasLimit(height)
	if gl_err != nil {
		return "", gl_err
	}
	r := strings.NewReader(archonAbi.Abi)
	scAbi, err_scAbi := abi.JSON(r) // reader io.Reader
	if err_scAbi != nil {
		return "", err_scAbi
	}

	bNodeID := []byte(params.NodeID)
	var b32NodeID [32]byte
	copy(b32NodeID[:], bNodeID[2:])

	args := encodings.ArchonSCArgs{Params: encodedParams,
		NodeID:        b32NodeID,
		HardwareProof: params.HardwareProof}
	var methodName string = "registerSP"
	dataFormatted, df_err := encodings.FormatData(scAbi, methodName, args)
	if df_err != nil {
		return "", df_err
	}
	contractAddressString := strings.Replace(archonAbi.ContractAddress(), "0x", "", 1)
	var contractAddress []byte
	for i := 0; i < len(contractAddressString); i += 2 {
		r, _ := hexutil.Decode("0x" + contractAddressString[i:i+2])
		contractAddress = append(contractAddress, []byte(r)...)
	}
	var bContractAddress [20]byte
	for i := 0; i < 20; i++ {
		bContractAddress[i] = contractAddress[i]
	}
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Wallet.Address, bContractAddress, amount, dataFormatted)
	if gp_err != nil {
		return "", gp_err
	}
	accountHasEnoughEthers, balance, totalCost, err := rpc_utils.CheckTxCostAgainstBalance(params.Stake, gasLimit, params.Wallet.Address)
	if err != nil {
		return "", err
	}
	if !accountHasEnoughEthers {
		return "", fmt.Errorf("error RegisterSP: totalCost of tx is ", totalCost.Text(10), " but account balance is ", balance.Text(10))
	}
	tx := types.NewTransaction(nonce,
		bContractAddress,
		amount,
		gasLimit,
		gasPrice,
		dataFormatted)
	signedTx, err_signedTx := params.Wallet.SignTx(tx, height)
	if err_signedTx != nil {
		return "", err_signedTx
	}

	// send tx
	txidString, txid_err := rpc_utils.SendRawTx(signedTx)
	if txid_err != nil {
		return "", txid_err
	}
	return txidString, nil
}

func UnregisterSP(params SPParams) (ret string, err error) {
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
	amount := new(big.Int)
	amount.SetUint64(uint64(0))
	height, h_err := rpc_utils.GetBlockHeight()
	if h_err != nil {
		return "", h_err
	}
	gasLimit, gl_err := rpc_utils.GetGasLimit(height)
	if gl_err != nil {
		return "", gl_err
	}
	r := strings.NewReader(archonAbi.Abi)
	scAbi, err_scAbi := abi.JSON(r) // reader io.Reader
	if err_scAbi != nil {
		return "", err_scAbi
	}

	args := encodings.ArchonSCArgs{}
	var methodName string = "unregisterSP"
	dataFormatted, df_err := encodings.FormatData(scAbi, methodName, args)
	if df_err != nil {
		return "", df_err
	}

	contractAddressString := strings.Replace(archonAbi.ContractAddress(), "0x", "", 1)
	var contractAddress []byte
	for i := 0; i < len(contractAddressString); i += 2 {
		r, _ := hexutil.Decode("0x" + contractAddressString[i:i+2])
		contractAddress = append(contractAddress, []byte(r)...)
	}
	var bContractAddress [20]byte
	for i := 0; i < 20; i++ {
		bContractAddress[i] = contractAddress[i]
	}
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Wallet.Address, bContractAddress, amount, dataFormatted)
	if gp_err != nil {
		return "", gp_err
	}
	accountHasEnoughEthers, balance, totalCost, err := rpc_utils.CheckTxCostAgainstBalance(uint64(0), gasLimit, params.Wallet.Address)
	if err != nil {
		return "", err
	}
	if !accountHasEnoughEthers {
		return "", fmt.Errorf("error RegisterSP: totalCost of tx is ", totalCost, " but account balance is ", balance)
	}
	tx := types.NewTransaction(nonce,
		bContractAddress,
		amount,
		gasLimit,
		gasPrice,
		dataFormatted)
	signedTx, err_signedTx := params.Wallet.SignTx(tx, height)
	if err_signedTx != nil {
		return "", err_signedTx
	}
	// send tx
	unregisterSPTxId, u_err := rpc_utils.SendRawTx(signedTx)
	if u_err != nil {
		return "", u_err
	}
	return unregisterSPTxId, nil
}
