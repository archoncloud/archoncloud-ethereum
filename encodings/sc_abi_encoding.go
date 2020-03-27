package encodings

import (
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ArchonSCArgs struct {
	Params               [32]byte
	NodeID               [32]byte
	HardwareProof        [32]byte
	Address              [20]byte
	HashedArchonFilepath [32]byte
	ContainerSignature   [ethcrypto.SignatureLength]byte
	Shardsize            uint64
	SPsToUploadTo        [][20]byte
	Username             [32]byte
	PublicKey            [64]byte
}

func FormatData(abi abi.ABI, methodName string, args ArchonSCArgs) (ret []byte, err error) {
	if methodName == "registerSP" {
		packedData, err_packedData := abi.Pack(methodName,
			args.Params,
			args.NodeID,
			args.HardwareProof) // ARGS
		if err_packedData != nil {
			return *new([]byte), err_packedData
		}
		return packedData, nil
	} else if methodName == "spAddress2SPProfile" {
		packedData, err_packedData := abi.Pack(methodName, args.Address) // ARGS
		if err_packedData != nil {
			return *new([]byte), err_packedData
		}
		return packedData, nil
	} else if methodName == "unregisterSP" {
		packedData, err_packedData := abi.Pack(methodName) // ARGS
		if err_packedData != nil {
			return *new([]byte), err_packedData
		}
		return packedData, nil
	} else if methodName == "proposeUpload" {
		var containerSignatureR, containerSignatureS [32]byte
		copy(containerSignatureR[:], args.ContainerSignature[0:32])
		copy(containerSignatureS[:], args.ContainerSignature[32:64])
		// note: V is stashed in param encoding
		packedData, err_packedData := abi.Pack(methodName,
			args.HashedArchonFilepath,
			containerSignatureR,
			containerSignatureS,
			args.Params,
			args.Shardsize,
			args.SPsToUploadTo) // ARGS
		if err_packedData != nil {
			return *new([]byte), err_packedData
		}
		return packedData, nil
	} else if methodName == "registerUsername" {
		var publicKeyX, publicKeyY [32]byte
		copy(publicKeyX[:], args.PublicKey[0:32])
		copy(publicKeyY[:], args.PublicKey[33:64])
		packedData, err_packedData := abi.Pack(methodName, args.Username, publicKeyX, publicKeyY) // ARGS
		if err_packedData != nil {
			return *new([]byte), err_packedData
		}
		return packedData, nil
	}
	var empty []byte
	return empty, fmt.Errorf("error FormatData, no data formatted")
}
