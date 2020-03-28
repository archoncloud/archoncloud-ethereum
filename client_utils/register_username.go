package client_utils

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
	"github.com/archoncloud/archoncloud-ethereum/encodings"
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
	"github.com/archoncloud/archoncloud-ethereum/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethparams "github.com/ethereum/go-ethereum/params"
)

type RegisterUsernameParams struct {
	Username string
	Wallet   wallet.EthereumKeyset
}

func RegisterUsername(params *RegisterUsernameParams) (txid string, err error) {
	var ret string = ""
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
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
	if len(params.Username) > 32 {
		return ret, fmt.Errorf("error RegisterUsername, username must be <= 32 chars")
	}
	var bUsername [32]byte
	copy(bUsername[:], []byte(params.Username)[0:32])

	args := encodings.ArchonSCArgs{Username: bUsername,
		PublicKey: params.Wallet.PublicKey}

	var methodName string = "registerUsername"
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
	copy(bContractAddress[0:20], contractAddress[0:20])
	amount := new(big.Int)
	amount.SetInt64(int64(0)) // trivial pmt
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Wallet.Address,
		bContractAddress,
		amount,
		dataFormatted)
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

	// sign tx
	point, err_point := ethcrypto.ToECDSA(params.Wallet.PrivateKey[:])
	if err_point != nil {
		return "", err_point
	}
	bigHeight := new(big.Int)
	baseHeight, err_res := strconv.ParseUint(strings.Replace(height, "0x", "", 1), 16, 64)
	if err_res != nil {
		return "", err_res
	}
	bigHeight.SetUint64(baseHeight)
	var signer types.Signer
	if archonAbi.RpcUrl() == "http://127.0.0.1:7545" || archonAbi.ChainIs() == "Ganorge" {
		signer = types.NewEIP155Signer(big.NewInt(g_chainID))
	} else {
		// gorli
		signer = types.MakeSigner(ethparams.GoerliChainConfig, bigHeight)
	}
	signedTx, err_signedTx := types.SignTx(tx, signer, point)
	if err_signedTx != nil {
		return "", err_signedTx
	}
	// send tx
	txid, tx_err := rpc_utils.SendRawTx(signedTx)
	if tx_err != nil {
		return "", tx_err
	}
	ret = txid
	return ret, nil
}
