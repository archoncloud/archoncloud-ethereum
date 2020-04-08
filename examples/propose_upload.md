# Example: Propose Upload

```
import (
	"math/rand" 
	// for example values 
	
	"github.com/archoncloud/archoncloud-ethereum/client_utils"
	"github.com/archoncloud/archoncloud-ethereum/wallet"

	// ...	
)

	// ...


		var keystoreFilepath string = "testingWallet.json"
		var password string = "TestingWallet"
		keyset, err := wallet.GetEthKeySet(keystoreFilepath, password)
		if err != nil {
			// handle
		}

		var serviceDuration = uint32(rand.Intn(24))
		// using rand for example values
		var minSLARequirements = rand.Intn(255)
		var containerSignature [65]byte
		for i := 0; i < 65; i++ {
			containerSignature[i] = byte(rand.Intn(256))
			// for example using random values. 
			// The container signature is a ECDSA 
			// signature of the upload contents
		}
		var archonFilepath string = "example/archonfilepath/a/b/c/d"
		var filesize uint64 = 654009 // example filesize
    		var shardsize uint64 = 6150 // example 
		
		fileContainerType := uint8(rand.Uint32()) 
		encryptionType := uint8(rand.Uint32())
		compressionType := uint8(rand.Uint32())
		shardContainerType := uint8(rand.Uint32())
		erasureCodeType := uint8(rand.Uint32())
		// see archoncloud-go for standard enum values for these
		// entries

		var addressStrings [5]string = [5]string{"0x595e356DDF600fea06a495731b739611b39e51E4",
			"0x0b22b2aB87646F23481e544358a257673385bdAa",
			"0xD826f70eD892D03D2CBc28427308bCB191993Bcc",
			"0xa4DC6292518245fE56145CF7Ee2660ba02376F20",
			"0x0990CA3842d2D45E5D488fe088FF4aE662e0bB5B"} 
			// example addresses determined by marketplace
		
		var spAddressesToUploadTo [][20]byte = ConvertAddressesToBytes(addressStrings) // pseudocode

    		uploadPmt := serviceDuration * maxMinAskPrice * len(addressStrings) 
		p := client_utils.UploadParams{
			Wallet: 	    keyset,
			ServiceDuration:    uint32(serviceDuration), //int
			MinSLARequirements: minSLARequirements,
			UploadPmt:          uint64(uploadPmt), // uint64
			ArchonFilepath:     archonFilepath, //string
			Filesize:           filesize,       //int
      			Shardsize:          shardsize,
      			FileContainerType:  fileContainerType,
			EncryptionType:     encryptionType,
			CompressionType:    compressionType,
			ShardContainerType: shardContainerType,
			ErasureCodeType:    erasureCodeType,
			ContainerSignature: containerSignature, //[32]byte
			SPsToUploadTo:      spAddressesToUploadTo}

		txid, err := client_utils.ProposeUpload(&p)
		if err != nil {
			// handle
		}
```
