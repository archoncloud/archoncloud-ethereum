# archoncloud-ethereum

### Contents:

  1. Overview

  2. High-Level protocol description

  3. Initialization

  4. APIs 

   - for the SP (storage provider)

   - for The Uploader

   - for The Downloader

   - for All Entities
     
  5. Wallet

  6. ABI

  7. Encodings

  8. RPC Utils

--------------------------------------------------------------------


### 1. Overview

This module serves as an api interface for applications interacting with the Archon Cloud (AC) either as agents of the cloud, or as beneficiaries. Specifically, this module is an api interface to interact with the Archon Ethereum Smart Contract (SC). The AC currently has smart contracts live in both the Ethereum and Neo networks that manage this protocol in analogous ways. The AC was designed to integrate smart contracts from other blockchains as needed. In this repository we focus on the api to the Archon Ethereum Smart Contract, but keep in mind the implementation is similar in Archon's other blockchain layers.


--------------------------------------------------------------------

### 2. High-Level protocol description

This is a very high-level description. Many details are glossed over in order to keep this brief. For a more detailed description, see the Archon White Paper or read the source code from our official repositories (including this one).

For this simple protocol description, we define the players in the Archon Cloud to be storage providers S, uploaders U, and downloaders D. The intent of these players are what you think they would be: the U want to make their content available, the D want to obtain the content of U, and S wants to earn cryptocurrency by acting as a conduit serving the needs of U and D.

To bootstrap this protocol, we start with the S. Any storage provider S must be registered with the Archon Ethereum Smart Contract as storage providers. This registration includes providing information about their storage capabilities, marketplace ask, routing information, as well as staking Ethereum token. An uploader U must also be registered with the SC. This registration includes establishing a namespace, and publishing the public key corresponding to their pseudo-identity.

We follow an upload u from U to its final target, the downloader D.

To keep this description very simple, we abstract away the details of erasure-encoding, etc and just follow u from U to S to D.

The U prepares u using some encoding and cryptographically signs u to get {u,sig(u), {other-metadata}}. Either now, or in the past, U has accumulated a subset of storage providers S_ = {S_1,S_2,...,S_n} (a subset of all storage providers in AC) from one or a few S. Locally, U runs the AC marketplace to determine the best S from S_ to accomodate u.

U concurrently makes a proposeUpload transaction pu_tx to SC with metadata of u and S that was determined by the marketplace and sends {u, sig(u), {other-metadata}}. S caches this upload and listens to SC for pu_tx to be confirmed by the blockchain. The proposeUpload transaction includes a payment to S for storing u, documents metadata of u including sig(u), and also validates the result of the marketplace. Assuming pu_tx is confirmed, S announces to the networking [overlay](https://github.com/archoncloud/archon-dht) of AC that it is storing {u, ...} and stores {u,...} for the period paid for by U in pu_tx. 

The downloader D knows of u from some other channel. Perhaps U advertised on, say, reddit that U uploaded u. D contacts some storage provider S' asking for the AC download url of u. Storage provider S' queries its networking [overlay](https://github.com/archoncloud/archon-dht) for the url(s) of any S holding u and returns these values to D. Downloader D downloads {u, sig(u), {other-metadata}} from S and retrieves the public key of U from the SC. D validates sig(u) with this public key and accepts u in the ideal case.

We will see below which API's each of the players call in order to participate in this protocol. Please keep in mind, this description glossed over some very important implementation details in order to be brief. For a more detailed protocol description, refer to the Archon Whitepaper, or inspect the source of the official repositories (including this one).


--------------------------------------------------------------------

### 3. Initialization

The developer must point their application at their preferred Ethereum RPC Url.

```
	err := rpc_utils.SetRpcUrl([]string{"<rpc-url1>", "<rpc-url2>", "<rpc-url3>"})
	if err != nil {
		// This means that no entered url is valid 
	}
```

--------------------------------------------------------------------

### 4. APIs 

--------------------------------------------------------------------

#### for the SP (storage provider)

`func RegisterSP(params SPParams) (txid string, err error)`

To participate in AC as a storage provider, S makes this function call to the Archon Ethereum Smart Contract with suitable "SPParams". 

`func UnregisterSP(params SPParams) (ret string, err error)`

A storage provider can unregister with AC by calling this function. Unregistering deletes from the AC the profile of the storage provider, and pays the S the sum of its earnings and stake. Once unregistered, the SP can no longer be assigned uploads in AC or earn from such uploads.
 
`func GetUploadTx(txHash [32]byte) (uploadTx UploadTx, err error)`

When an upload u from uploader U is made to storage provider S in AC, the S obtains from the blockchain the upload transaction associated with u using txHash collated with u. The S compares this returned data with the upload u to ensure that the upload matches the parameters and payment made to the SC. The S polls the blockchain with this txid to be sure this transaction is confirmed before storing this upload in fulfillment of the Service Level Agreement.

`func GetUsernameFromContract(address [20]byte) (username [32]byte, err error)`

A subroutine of the previous function call, is getting the username from the contract that is registered to the uploader. This is a security precaution to prevent the uploader from maliciously overwriting files in other uploader's namespaces. The storage provider stores the upload u in the namespace returned from this function call. This does not necessarily need to be called directly by the developer but may be handy.

`func GetRegisteredSP(ethAddress [20]byte) (sp *RegisteredSp, err error)`

A courtesy that an S provides to the AC, is that it serves as a proxy to the AC to "light-clients". Uploaders and Downloaders can be light-clients. A way an S acts as a proxy is that S serves a cache of storage provider profiles corresponding to a census of it's known nodes in the network [overlay](https://github.com/archoncloud/archon-dht). The way the S forms and maintains this cache is by collecting storage provider data of these known nodes that is stored both in the SC and in the network [overlay](https://github.com/archoncloud/archon-dht). The data that is stored in the SC is retrieved using this "GetRegisteredSP" function.

`func GetNodeID2Address(nodeID [32]byte) ([20]byte, error)`

A subroutine of the caching process mentioned in the description of "GetRegisteredSP" is satisfied using this function. Each node in the networking [overlay](https://github.com/archoncloud/archon-dht) has a unique nodeID. So for each node known by S, the nodeID acts as a key or handle to the storage provider profile data stored in both the SC and the network [overlay](https://github.com/archoncloud/archon-dht). As far as retrieving storage provider profile data from the SC is concerned, the flow looks like nodeID -> ethaddress -> registeredSp (Profile). This does not necessarily need to be called directly by the developer, but is used in the networking [overlay](https://github.com/archoncloud/archon-dht).

--------------------------------------------------------------------

#### for the Uploader

`func RegisterUsername(params *RegisterUsernameParams) (txid string, err error)` 

A necessary condition of an upload in AC being considered valid, is that it's uploader U is registered with the SC. This registration establishes in the blockchain storage the namespace and the pubic key of U.

`func ProposeUpload(params *UploadParams) (txid string, err error)` 

An uploader calls this function with appropriate parameters as a step in the upload process. This function constructs, signs, and broadcasts to the blockchain an upload transaction call on the SC. This transaction includes a list of storage providers who are to receive the upload (shards), payment, and file integrity data such as container signature. The list of storage providers and payment is settled in the local marketplace instance handled in other modules, [for marketplace see github.com/archoncloud/archoncloud-go](https://github.com/archoncloud/archoncloud-go)

--------------------------------------------------------------------

#### for the downloader

`func GetPublickeyFromContract(username string, timeout time.Duration) (pubkey [64]byte, err error)`

Given an upload u from uploader U into AC under it's namespace "username", a downloader who downloads u want to validate it's integrity. The downloader validates the cryptographic signature on u using the public key of U. This public key is obtained by calling this function.


--------------------------------------------------------------------

#### for all entities

`func GetBalance(ethAddress [20]byte) (big.Int, error)`

All Ethereum Addresses have a non-negative integer balance in wei.


`func GetEarnings(ethAddress [20]byte) (big.Int, error)`

The players in AC who acrue earnings are the storage providers. Any entity with access to the SC can view this public information.


`func GetTxLogs(txid string) (TxLogs, error)`

Many Ethereum Smart Contract functions emit logs. This function is useful to access logs corresponding to a given transaction id.

--------------------------------------------------------------------

####  5. Wallet

`func GetEthKeySet(keystoreFilepath, password string) (ethKeySet EthereumKeyset, err error)`

Given keystoreFilepath and password, open an encrypted V3 keystore file to obtain EthereumKeyset.

`func (e *EthereumKeyset) SignTx(tx *types.Transaction, height string) (*types.Transaction, error)`

This SignTx method is used in the repository as a subroutine of the API's that construct, sign, and broadcast transactions to the Ethereum Blockchain. The developer does not need to call this directly in AC.

`func (e *EthereumKeyset) ExportPrivateKey() (string, error)`

CAUTION! Exposing your private key may lead to a loss of Ethereum token! Do not use this function unless you are familiar with best security practices with respect to public-key cryptography. 

`func GenerateAndSaveEthereumWallet(keystoreFilepath, password string) error`

Generate encrypted V3 keystore with chosen password and store to filesystem at keystoreFilepath.

`func GenerateAndSaveEthereumWalletFromPrivateKey(privateKey, keystoreFilepath, password string) error`

Given privateKey, generate encrypted V3 keystore with chosen password and store to filesystem at keystoreFilepath.

--------------------------------------------------------------------

### 6. ABI

The ABI (application binary interface) along with the Archon Ethereum Smart Contract address are kept in the `abi/abi.go` file along with other minor utilities. For the most part, this file has no API's that would be useful to the developer and should be left alone. The AC protocol relies on the Archon Ethereum Smart Contract as a control to various parts of the protocol. Said differently, each player in the AC protocol acts in response to the evolving state of the Archon Ethereum Smart Contract in conjunction with other events. By design, it is not advantageous in the AC protocol for an entity to manipulate its interface with the SC, or the SC itself (changing address in this file). The design effectively excises such entities from the protocol. "So don't bother".

--------------------------------------------------------------------

### 7. Encodings

Like the ABI file above, the `encodings` folder contains assets that are necessary to this module, but are not intended to be used directly by the developer. Many of the API's exposed by this repository require various encodings/decodings behind the scenes. 

--------------------------------------------------------------------

### 8. RPC Utils

Again, the `rpc_utils` folder contains functions that are not intended to be of immediate use to the developers using this repository. Rather the functions provided by the `rpc_utils` folder are called as subroutines of the showcased API's provided by this repository.

--------------------------------------------------------------------
