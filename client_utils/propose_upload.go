package client_utils

import (
	"fmt"
	"math/big"

	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"
	"github.com/archoncloud/archoncloud-ethereum/encodings"
	"github.com/archoncloud/archoncloud-ethereum/rpc_utils"
	"github.com/archoncloud/archoncloud-ethereum/wallet"
)

var g_chainID int64 = archonAbi.ChainID()

type UploadParams struct {
	Wallet wallet.EthereumKeyset

	ServiceDuration    uint32
	MinSLARequirements int
	UploadPmt          uint64 // bid in marketplace, defaults to flat payment
	ArchonFilepath     string
	Filesize           uint64
	Shardsize          uint64
	FileContainerType  uint8
	EncryptionType     uint8
	CompressionType    uint8
	ShardContainerType uint8
	ErasureCodeType    uint8
	CustomField        uint8

	ContainerSignature [ethcrypto.SignatureLength]byte
	SPsToUploadTo      [][20]byte // sps to whom the shards are to go
}

// To upload shards to sps, the uploader must make ProposeUpload tx to the
// smart contract with "SPsToUploadTo" being the result of the local
// marketplace instance with the upload as input. When the uploader sends
// the shards to the sps with the resultant txid, the sps will check that the
// txid has has data that matches the upload.. i.e. correct upload size,
// containerSignature, etc.
func ProposeUpload(params *UploadParams) (txid string, err error) {
	// construct tx
	nonce, nonce_err := rpc_utils.GetNonceForAddress(params.Wallet.Address)
	if nonce_err != nil {
		return "", nonce_err
	}
	amount := new(big.Int)
	amount.SetUint64(params.UploadPmt)
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
	bArchonFilepath := []byte(params.ArchonFilepath)
	hashedArchonFilepath := ethcrypto.Keccak256(bArchonFilepath[:])
	proposeUploadParams := encodings.ProposeUploadParams{
		ServiceDuration:     params.ServiceDuration,
		MinSLARequirements:  params.MinSLARequirements,
		UploadPmt:           params.UploadPmt,
		Filesize:            params.Filesize,
		FileContainerType:   params.FileContainerType,
		EncryptionType:      params.EncryptionType,
		CompressionType:     params.CompressionType,
		ShardContainerType:  params.ShardContainerType,
		ErasureCodeType:     params.ErasureCodeType,
		ContainerSignatureV: params.ContainerSignature[64], // stashing V
		CustomField:         params.CustomField}
	encodedParams, ep_err := encodings.EncodeProposeUploadParams(
		proposeUploadParams)
	if ep_err != nil {
		return "", ep_err
	}
	var bHashedArchonFilepath [32]byte
	copy(bHashedArchonFilepath[:], hashedArchonFilepath[0:32])
	args := encodings.ArchonSCArgs{HashedArchonFilepath: bHashedArchonFilepath,
		ContainerSignature: params.ContainerSignature,
		Params:             encodedParams,
		Shardsize:          params.Shardsize,
		SPsToUploadTo:      params.SPsToUploadTo}
	var methodName string = "proposeUpload"
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
	gasPrice, gp_err := rpc_utils.EstimateGas(params.Address,
		bContractAddress,
		amount,
		dataFormatted)
	if gp_err != nil {
		return "", gp_err
	}
	accountHasEnoughEthers, balance, totalCost, err := rpc_utils.CheckTxCostAgainstBalance(params.UploadPmt, gasLimit, params.Address)
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
	signedTx, err_signedTx := params.Wallet.SignTx(tx)
	if err_signedTx != nil {
		return "", err_signedTx
	}
	// send tx
	txid, tx_err := rpc_utils.SendRawTx(signedTx)
	if tx_err != nil {
		return "", tx_err
	}
	return txid, nil
}
