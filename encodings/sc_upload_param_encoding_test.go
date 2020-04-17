package encodings

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncodeProposeUpload(t *testing.T) {
	{
		tm := time.Now()
		tUnix := tm.Unix()
		rand.Seed(tUnix)
		var uploadPmt = uint64(rand.Intn(1000000))
		var serviceDuration = uint32(rand.Intn(3 * 256))
		var minSLARequirements = rand.Intn(17)
		var filesize uint64 = 654009
		fileContainerType := uint8(rand.Uint32())
		encryptionType := uint8(rand.Uint32())
		compressionType := uint8(rand.Uint32())
		shardContainerType := uint8(rand.Uint32())
		erasureCodeType := uint8(rand.Uint32())
		accessControlLevel := uint8(rand.Uint32())
		customField := uint8(rand.Uint32())

		p := ProposeUploadParams{
			ServiceDuration:    uint32(serviceDuration), //int
			UploadPmt:          uploadPmt,
			MinSLARequirements: minSLARequirements,
			Filesize:           filesize, //int
			FileContainerType:  fileContainerType,
			EncryptionType:     encryptionType,
			CompressionType:    compressionType,
			ShardContainerType: shardContainerType,
			ErasureCodeType:    erasureCodeType,
			AccessControlLevel: accessControlLevel,
			CustomField:        customField}

		enc, err_1 := EncodeProposeUploadParams(p)
		if err_1 != nil {
			assert.Equal(t, true, false, "encode Url error")
		}
		dec := DecodeProposeUploadParams(enc)

		assert.Equal(t, serviceDuration, dec.ServiceDuration, "encode,decode mismatch: serviceDuration")
		assert.Equal(t, uploadPmt, dec.UploadPmt, "encode,decode mismatch: uploadPmt")
		assert.Equal(t, filesize, dec.Filesize, "encode,decode mismatch: filesize")
		assert.Equal(t, accessControlLevel, dec.AccessControlLevel, "encode,decode mismatch: AccessControlLevel")
	}
}
