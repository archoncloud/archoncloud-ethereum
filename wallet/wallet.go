package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types" //NewTransaction
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethparams "github.com/ethereum/go-ethereum/params"

	archonAbi "github.com/archoncloud/archoncloud-ethereum/abi"

	"github.com/pborman/uuid"
)

type SafeFilename struct {
	sync.RWMutex
}

var Filenames = struct {
	sync.RWMutex
	M map[string]*SafeFilename
}{M: make(map[string]*SafeFilename)}

type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

type EthereumKeyset struct {
	privateKey [32]byte
	PublicKey  [64]byte
	Address    [20]byte
}

var g_chainID int64 = archonAbi.ChainID()

// Costfactor, CPU/memory cost parameter (must be power of 2)
var g_scryptN int = 8192

// Parallelization parameter
var g_scryptP int = 1

func GetEthKeySet(keystoreFilepath, password string) (ethKeySet EthereumKeyset, err error) {
	dat, err := ioutil.ReadFile(keystoreFilepath)
	if err != nil {
		return *new(EthereumKeyset), fmt.Errorf("GetEthKeySet error, file read error")
	}
	var w encryptedKeyJSONV3
	err_w := json.Unmarshal(dat, &w)
	if err_w != nil {
		return *new(EthereumKeyset), fmt.Errorf("GetEthKeySet error, json parse error. Check that wallet is in V3 keystore format.")
	}
	privateKey, err_decrypted := keystore.DecryptDataV3(w.Crypto, password)
	if err_decrypted != nil {
		return *new(EthereumKeyset), fmt.Errorf("invalid credentials")
	}
	point, err_point := ethcrypto.ToECDSA(privateKey)
	if err_point != nil {
		return *new(EthereumKeyset), fmt.Errorf("GetEthKeySet error, ecdsa error")
	}
	publicKeyStringX := hexutil.EncodeBig(point.PublicKey.X)
	publicKeyStringY := hexutil.EncodeBig(point.PublicKey.Y)
	for len(publicKeyStringX) < 66 {
		publicKeyStringX = "0x" + "0" + strings.Replace(publicKeyStringX, "0x", "", -1)
	}
	for len(publicKeyStringY) < 66 {
		publicKeyStringY = "0x" + "0" + strings.Replace(publicKeyStringY, "0x", "", -1)
	}
	publicKey := hexutil.MustDecode(publicKeyStringX + strings.Replace(publicKeyStringY, "0x", "", -1))
	var address []byte
	for i := 0; i < len(w.Address); i += 2 {
		r, _ := hexutil.Decode("0x" + w.Address[i:i+2])
		address = append(address, []byte(r)...)
	}
	var setPrivateKey [32]byte
	copy(setPrivateKey[:], privateKey[0:32])
	var setPublicKey [64]byte
	copy(setPublicKey[:], publicKey[0:64])
	var setAddress [20]byte
	copy(setAddress[:], address[0:20])
	ethereumKeyset := EthereumKeyset{privateKey: setPrivateKey, PublicKey: setPublicKey, Address: setAddress}
	return ethereumKeyset, nil
}

func (e *EthereumKeyset) SignTx(tx *types.Transaction, height string) (*types.Transaction, error) {
	point, err_point := ethcrypto.ToECDSA(e.privateKey[:])
	if err_point != nil {
		return tx, err_point
	}
	bigHeight := new(big.Int)
	baseHeight, err_res := strconv.ParseUint(strings.Replace(height, "0x", "", 1), 16, 64)
	if err_res != nil {
		return tx, err_res
	}
	bigHeight.SetUint64(baseHeight)
	var signer types.Signer
	if archonAbi.Rpc().Url == "http://127.0.0.1:7545" || archonAbi.ChainIs() == "Ganorge" {
		signer = types.NewEIP155Signer(big.NewInt(g_chainID))
	} else {
		// gorli
		signer = types.MakeSigner(ethparams.GoerliChainConfig, bigHeight)
	}
	signedTx, err_signedTx := types.SignTx(tx, signer, point)
	if err_signedTx != nil {
		return tx, err_signedTx
	}
	return signedTx, nil
}

func GenerateAndSaveEthereumWallet(keystoreFilepath, password string) error {
	key, err := generateWallet(rand.Reader)
	if err != nil {
		return err
	}
	scryptN := g_scryptN
	scryptP := g_scryptP
	encryptedKey, err := keystore.EncryptKey(key, password, scryptN, scryptP)
	if err != nil {
		return err
	}
	safeKeystoreFilepath, err := overwriteproofFilename(keystoreFilepath)
	if err != nil {
		return err
	}
	Filenames.Lock()
	if Filenames.M[safeKeystoreFilepath] == nil {
		Filenames.M[safeKeystoreFilepath] = new(SafeFilename)
	}
	ptr := Filenames.M[safeKeystoreFilepath]
	Filenames.Unlock()
	ptr.Lock()
	defer ptr.Unlock()
	ioutil.WriteFile(safeKeystoreFilepath, encryptedKey, 0644)
	return nil
}

func GenerateAndSaveEthereumWalletFromPrivateKey(privateKey, keystoreFilepath, password string) error {
	bigIntPrivKey := new(big.Int)
	bigIntPrivKey.SetString(privateKey[2:], 16)

	point, err := ethcrypto.ToECDSA(bigIntPrivKey.Bytes())
	if err != nil {
		return err
	}
	key := newKeyFromECDSA(point)

	scryptN := g_scryptN
	scryptP := g_scryptP
	encryptedKey, err := keystore.EncryptKey(key, password, scryptN, scryptP)
	if err != nil {
		return err
	}
	safeKeystoreFilepath, err := overwriteproofFilename(keystoreFilepath)
	if err != nil {
		return err
	}
	Filenames.Lock()
	if Filenames.M[safeKeystoreFilepath] == nil {
		Filenames.M[safeKeystoreFilepath] = new(SafeFilename)
	}
	ptr := Filenames.M[safeKeystoreFilepath]
	Filenames.Unlock()
	ptr.Lock()
	defer ptr.Unlock()
	ioutil.WriteFile(safeKeystoreFilepath, encryptedKey, 0644)
	return nil
}

//CAUTION! Exposing your private key may lead to a loss of Ethereum token! Do not use this function unless you are familiar with best security practices with respect to public-key cryptography.
func (e *EthereumKeyset) ExportPrivateKey() (string, error) {
	var bPrivateKey []byte
	for i := 0; i < len(e.privateKey); i++ {
		bPrivateKey = append(bPrivateKey, e.privateKey[i])
	}
	ret := hexutil.Encode(bPrivateKey)
	return ret, nil
}

//CAUTION! Exposing your private key may lead to a loss of Ethereum token! Do not use this function unless you are familiar with best security practices with respect to public-key cryptography.
func (e *EthereumKeyset) ExportPrivateKeyBytes() ([]byte, error) {
	var bPrivateKey []byte
	for i := 0; i < len(e.privateKey); i++ {
		bPrivateKey = append(bPrivateKey, e.privateKey[i])
	}
	return bPrivateKey
}

func generateWallet(rand io.Reader) (*keystore.Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(ethcrypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	key := newKeyFromECDSA(privateKeyECDSA)
	if !strings.HasPrefix(key.Address.Hex(), "0x00") {
		return generateWallet(rand)
	}
	return key, nil
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *keystore.Key {
	// function "lifted" from go-ethereum
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

// If path/to/wallet/exampleWallet.json already exists on your system,
// this function will output path/to/wallet/exampleWallet(1).json so that
// the first wallet will not be overwritten
// This generalizes to path/to/wallet/exampleWallet(n).json ->
//   path/to/wallet/exampleWallet(n+1).json for n in N
func overwriteproofFilename(filepath string) (string, error) {
	Filenames.Lock()
	if Filenames.M[filepath] == nil {
		Filenames.M[filepath] = new(SafeFilename)
	}
	ptr := Filenames.M[filepath]
	Filenames.Unlock()
	ptr.Lock()
	defer ptr.Unlock()
	_, err := os.Stat(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		} else {
			// os.IsNotExist(err)
			// do nothing
			return filepath, nil
		}
	} else {
		newFilepath := renameFilenameWithDistinctIndex(filepath)
		return overwriteproofFilename(newFilepath)
	}
	return "", nil
}

func renameFilenameWithDistinctIndex(filepath string) string {
	newFilepath := ""
	// check if has "."
	periodIndex := strings.Index(filepath, ".")
	if periodIndex > -1 {
		// GREP FOR SUFFIX
		suffixRegex := regexp.MustCompile("\\..+$")
		wOutSuffix := suffixRegex.ReplaceAllString(filepath, "")
		preSuffixRegex := regexp.MustCompile("^.*\\.")
		suffix := preSuffixRegex.ReplaceAllString(filepath, "")
		// GREP FOR INDEX
		indexRegex := regexp.MustCompile("\\([0-9]\\)$")
		wOutIndex := indexRegex.ReplaceAllString(wOutSuffix, "")

		preIndexRegex := regexp.MustCompile("^.*\\(")
		index := preIndexRegex.ReplaceAllString(wOutSuffix, "")
		index = strings.Replace(index, ")", "", 1)
		// INCREMENT INDEX
		indexUint, _ := strconv.ParseUint(index, 10, 10)
		indexUint++
		index = "(" + strconv.Itoa(int(indexUint)) + ")"

		newFilepath = wOutIndex + index + "." + suffix
	} else {
		preIndexRegex := regexp.MustCompile("^.*\\(")
		index := preIndexRegex.ReplaceAllString(filepath, "")
		index = strings.Replace(index, ")", "", 1)
		// INCREMENT INDEX
		indexUint, _ := strconv.ParseUint(index, 10, 10)
		indexUint++
		index = "(" + strconv.Itoa(int(indexUint)) + ")"

		postIndexRegex := regexp.MustCompile("\\(.*$")
		woutIndex := postIndexRegex.ReplaceAllString(filepath, "")
		newFilepath = woutIndex + index
	}
	return newFilepath
}
