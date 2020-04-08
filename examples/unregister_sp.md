# Example Unregister SP

```
import (
	"github.com/archoncloud/archoncloud-ethereum/register"
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
	p := register.SPParams{Wallet: keyset}
	txid, err := register.UnregisterSP(p)
	if err != nil {
		// handle
	}
```
