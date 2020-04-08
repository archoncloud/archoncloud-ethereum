# Example: Register SP

```
import (
	"math/rand" 
	// for example values 
	
	"github.com/archoncloud/archoncloud-ethereum/register"
	"github.com/archoncloud/archoncloud-ethereum/wallet"

	"github.com/pariz/gountries"

	// ...	
)

	// ...

	var keystoreFilepath string = "testingWallet.json"
	var password string = "TestingWallet"
	keyset, err := wallet.GetEthKeySet(keystoreFilepath, password)
	if err != nil {
		// handle
	}

	// using rand values for example
	var slaLevel = rand.Intn(8) // must be in range [1,8]
	var pledgedStorage = rand.Intn(5300000000)
	var bandwidth = rand.Intn(5300000000)
	JP := "JP" // Japan for example
	countryCode := gountries.Codes{Alpha2: JP}
	var minAskPrice = rand.Intn(3500000)
	var stake = 1000000000000000 // in wei. Equal to 0.01 Eth
	// current minimum stake. See Archon Smart Contract for stake as function of proposed utility 
	var hardwareProof [32]byte 
	// hardware proof in development. See Archon Whitepaper
	for i := 0; i < 32; i++ {
		hardwareProof[i] = byte(rand.Intn(256)) 
		// fake values for example
	}

	nodeID := overlay.GetNodeID() 
	// default archon network overlay is archon-dht
	p := register.SPParams{
			Wallet: keyset,
			SLALevel:         slaLevel,
			PledgedStorage:   uint64(pledgedStorage), 
			Bandwidth:        uint64(bandwidth),        
			CountryCode:      countryCode,
			MinAskPrice:      uint64(minAskPrice),
			Stake:            uint64(stake),
			HardwareProof:    hardwareProof,
			NodeID:           nodeID} 
    	txid, err := register.RegisterSP(p)
	if err != nil {
		// handle
	}

```
