package encodings

import (
	"fmt"
)

type ProposeUploadParams struct {
	ServiceDuration     uint32
	MinSLARequirements  int
	UploadPmt           uint64
	Filesize            uint64
	FileContainerType   uint8
	EncryptionType      uint8
	CompressionType     uint8
	ShardContainerType  uint8
	ErasureCodeType     uint8
	AccessControlLevel  uint8
	ContainerSignatureV byte
	CustomField         uint8
}

func EncodeProposeUploadParams(params ProposeUploadParams) (res [32]byte, err error) {
	var ret [32]byte
	// service duration
	// uint32 > 4bytes
	ret[0] = byte((params.ServiceDuration & 0x00000000FF000000) >> 24)
	ret[1] = byte((params.ServiceDuration & 0x0000000000FF0000) >> 16)
	ret[2] = byte((params.ServiceDuration & 0x000000000000FF00) >> 8)
	ret[3] = byte((params.ServiceDuration & 0x00000000000000FF) >> 0)
	// minSLARequirements
	// byte
	if params.MinSLARequirements > 255 {
		return ret, fmt.Errorf("MinSLARequirements out of range")
	}
	ret[4] = byte(params.MinSLARequirements)
	// uploadPmt
	// uint64 > 8bytes //
	ret[5] = byte((params.UploadPmt & 0xFF00000000000000) >> 56)
	ret[6] = byte((params.UploadPmt & 0x00FF000000000000) >> 48)
	ret[7] = byte((params.UploadPmt & 0x0000FF0000000000) >> 40)
	ret[8] = byte((params.UploadPmt & 0x000000FF00000000) >> 32)
	ret[9] = byte((params.UploadPmt & 0x00000000FF000000) >> 24)
	ret[10] = byte((params.UploadPmt & 0x0000000000FF0000) >> 16)
	ret[11] = byte((params.UploadPmt & 0x000000000000FF00) >> 8)
	ret[12] = byte((params.UploadPmt & 0x00000000000000FF) >> 0)
	// filesize
	// uint64 > 8bytes
	ret[13] = byte((params.Filesize & 0xFF00000000000000) >> 56)
	ret[14] = byte((params.Filesize & 0x00FF000000000000) >> 48)
	ret[15] = byte((params.Filesize & 0x0000FF0000000000) >> 40)
	ret[16] = byte((params.Filesize & 0x000000FF00000000) >> 32)
	ret[17] = byte((params.Filesize & 0x00000000FF000000) >> 24)
	ret[18] = byte((params.Filesize & 0x0000000000FF0000) >> 16)
	ret[19] = byte((params.Filesize & 0x000000000000FF00) >> 8)
	ret[20] = byte((params.Filesize & 0x00000000000000FF) >> 0)
	//
	ret[21] = byte(params.FileContainerType)
	ret[22] = byte(params.EncryptionType)
	ret[23] = byte(params.CompressionType)
	ret[24] = byte(params.ShardContainerType)
	ret[25] = byte(params.ErasureCodeType)
	ret[26] = byte(params.AccessControlLevel)

	// V for containerSignature stashed here
	ret[27] = params.ContainerSignatureV
	//
	ret[28] = byte(params.CustomField)

	// NOTE THERE ARE MANY FREE BYTES

	return ret, nil
}

func DecodeProposeUploadParams(params [32]byte) ProposeUploadParams {
	var ret ProposeUploadParams
	// service duration
	// uint32 > 4bytes
	var serviceDuration uint32 = 0
	for i := 0; i < 4; i++ {
		serviceDuration += uint32(params[i]) << uint((3-i)*8)
	}
	ret.ServiceDuration = serviceDuration
	// minSLARequirements
	// byte
	var minSLARequirements int = int(params[4])
	ret.MinSLARequirements = minSLARequirements
	// uploadPmt
	// uint64 > 8bytes
	var uploadPmt uint64 = 0
	for i := 0; i < 8; i++ {
		uploadPmt += uint64(uint64(params[i+5]) << uint((7-i)*8))
	}
	ret.UploadPmt = uploadPmt
	// filesize
	// uint64 > 8bytes
	var filesize uint64 = 0
	for i := 0; i < 8; i++ {
		filesize += uint64(uint64(params[i+13]) << uint((7-i)*8))
	}
	ret.Filesize = filesize
	ret.FileContainerType = uint8(params[21])
	ret.EncryptionType = uint8(params[22])
	ret.CompressionType = uint8(params[23])
	ret.ShardContainerType = uint8(params[24])
	ret.ErasureCodeType = uint8(params[25])
	ret.AccessControlLevel = uint8(params[26])
	// ContainerSignature V stashed here
	ret.ContainerSignatureV = byte(params[27])
	ret.CustomField = uint8(params[28])

	return ret
}
