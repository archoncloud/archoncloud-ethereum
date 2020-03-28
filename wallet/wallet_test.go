package wallet

import (
	//  "fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
) // import

func TestGetEthSet(t *testing.T) {
	var keystoreFilepath string = "testing_wallet_1.json"
	var password string = "ganacheTestingWallet1"
	keyset, err := GetEthKeySet(keystoreFilepath, password)
	if err != nil {
		assert.Equal(t, true, false, "error GetEthKeySet")
	}
	// tests
	var publicKey []byte
	publicKey = append(publicKey, keyset.PublicKey[:]...)
	var address []byte
	address = append(address, keyset.Address[:]...)

	check2 := hexutil.Encode(publicKey)
	check3 := hexutil.Encode(address)

	var expectedCheck2 string = "0xbb101d5d077c3ef1dc65b94811ca86935243d875567c05173269a6709680805861022ab167e340a51b340f695457f99eddc99caedaaf9382e62fb0653825be85"
	var expectedCheck3 string = "0x595e356ddf600fea06a495731b739611b39e51e4"
	assert.Equal(t, expectedCheck2, check2, "public key incorrect")
	assert.Equal(t, expectedCheck3, check3, "address incorrect")
}

func TestRenameFilenameWithDistinctIndex(t *testing.T) {
	//func renameFilenameWithDistinctIndex(filepath string) string

	initial := "test.json"
	result := renameFilenameWithDistinctIndex(initial)
	expected := "test(1).json"
	assert.Equal(t, expected, result, "incorrect file renaming")
	initial = expected
	expected = "test(2).json"
	result = renameFilenameWithDistinctIndex(initial)
	assert.Equal(t, expected, result, "incorrect file renaming")
	withoutSuffix := "testyboi"
	initial = withoutSuffix
	result = renameFilenameWithDistinctIndex(initial)
	expected = withoutSuffix + "(1)"
	assert.Equal(t, expected, result, "incorrect file renaming")
}
