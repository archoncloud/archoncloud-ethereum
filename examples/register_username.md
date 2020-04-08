# Example: Register Username

```
import (
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

	regParams := client_utils.RegisterUsernameParams{
		Username: "eXampLeUserNaME",
		Wallet:   keyset}
	txid, err := client_utils.RegisterUsername(&regParams)
	if err != nil {
		// handle
	}
```
