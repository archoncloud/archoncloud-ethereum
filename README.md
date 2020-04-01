# archoncloud-ethereum

NOTE: This software is in development and subject to change. Not for production use!

### Contents:

 - Overview

 - High-Level protocol description

 - APIs 

   - for the SP (storage provider)

   - for The Uploader

   - for The Downloader

   - for All Entities

     - Basic Getters
     
  - Wallet

  - ABI

  - Encodings

--------------------------------------------------------------------


### Overview

This module serves as an api interface for applications interacting with the Archon Decentralized Cloud (ADC) either as agents of the cloud, or as beneficiaries. Specifically, this module is an api interface to interact with the Archon Ethereum Smart Contract. The ADC currently has smart contracts live in both the Ethereum and Neo networks that manage this protocol in analogous ways. The ADC was designed to integrate smart contracts from other blockchains as needed. In this repository we focus on the api to the Archon Ethereum Smart Contract, but keep in mind the implementation is similar in Archon's other blockchain layers.


--------------------------------------------------------------------

### High-Level protocol description

This is a very high-level description. Many details are glossed over in order to keep this brief. For a more detailed description, see the Archon White Paper /TODO LINK/ or read the source code from our official repositories.

For this simple protocol description, we define the players in the Archon Decentralized Cloud to be storage providers S, uploaders U, and downloaders D. The intent of these players are what you think they would be: the U want to make their content available, the D want to obtain the content of U, and S wants to earn cryptocurrency by acting as a conduit serving the needs of U and D.

To bootstrap this protocol, we start with the S. Any storage provider S must be registered with the Archon Ethereum Smart Contract as storage providers. This registrations includes providing information about their storage capabilities, marketplace ask, routing information, as well as staking Ethereum token. An uploader U must also be registered with the SC. This registration includes establishing a namespace, and publishing the public key corresponding to their pseudo-identity.

We follow an upload u from U to its final target, the downloader D.

To keep this description very simple, we abstract away the details of erasure-encoding, etc and just follow u from U to S to D.

The U prepares u using some encoding and cryptographically signs u to get {u,sig(u), {other-metadata}}. Either now, or in the past, U has accumulated a subset of storage providers S_ = {S_1,S_2,...,S_n} (a subset of all storage providers in ADC) from one or a few S. Locally, U runs the ADC marketplace to determine the best S from S_ to accomodate u.

U concurrently makes a proposeUpload transaction pu_tx to SC with metadata of u and S that was determined by the marketplace and sends {u, sig(u), {other-metadata}}. S caches this upload and listens to SC for pu_tx to be confirmed by the blockchain. The proposeUpload transaction includes a payment to S for storing u, documents metadata of u including sig(u), and also validates the result of the marketplace. Assuming pu_tx is confirmed, S announces to the networking overlay of ADC that it is storing {u, ...} and stores {u,...} for the period paid for by U in pu_tx. 

The downloader D knows of u from some other channel. Perhaps U advertised on, say, reddit that U uploaded u. D contacts some storage provider S' asking for the ADC download url of u. Storage provider S' queries its networking overlay for the url(s) of any S holding u and returns these values to D. Downloader D downloads {u, sig(u), {other-metadata}} from S and retrieves the public key of U from the SC. D validates sig(u) with this public key and accepts u in the ideal case.

We will see below which API's each of the players call in order to participate in this protocol. Please keep in mind, this description glossed over some very important implementation details in order to be brief. For a more detailed protocol description, refer to the Archon Whitepaper /TODO NEED URL/, or inspect the source of the official repositories.

--------------------------------------------------------------------

### APIs 

--------------------------------------------------------------------

#### for the SP (storage provider)

###### functions

`func RegisterSP(params SPParams) (txid string, err error)`

To participate in ADC as a storage provider, S makes this function call to the Archon Ethereum Smart Contract with suitable "SPParams". 

`func UnregisterSP(params SPParams) (ret string, err error)`

A storage provider can unregister with ADC by calling this function. Once unregistered, the SP can no longer be assigned uploads in ADC or earn from such uploads. The storage provider's profile remains with the SC including value of earnings and remaining stake. The storage provider can withdrawal this token from the SC by calling ArchonSPWithdrawal method on the SC. (wrapper coming soon)

`func GetUploadTx(txHash [32]byte) (uploadTx UploadTx, err error)`

When an upload u from uploader U is made to storage provider S in ADC, the S obtains from the blockchain the upload transaction associated with u using txHash collated with u. The S compares this returned data with the upload u to ensure that the upload matches the parameters and payment made to the SC. The S polls the blockchain with this txid to be sure this transaction is confirmed before storing this upload in fulfillment of the Service Level Agreement.

`func GetUsernameFromContract(address [20]byte) (username [32]byte, err error)`


`func GetRegisteredSP(ethAddress [20]byte) (sp *RegisteredSp, err error)`


`func GetNodeID2Address(nodeID [32]byte) ([20]byte, error)`

--------------------------------------------------------------------

#### for the Uploader

###### functions

`func RegisterUsername(params *RegisterUsernameParams) (txid string, err error)` 

`func ProposeUpload(params *UploadParams) (txid string, err error)` 


--------------------------------------------------------------------

#### for the downloader

###### functions

`func GetPublickeyFromContract(username string, timeout time.Duration) (pubkey [64]byte, err error)`


--------------------------------------------------------------------

#### for all entities

###### functions

`func GetBalance(ethAddress [20]byte) (big.Int, error)`


`func GetEarnings(ethAddress [20]byte) (big.Int, error)`


`func GetTxLogs(txid string) (TxLogs, error)`


--------------------------------------------------------------------

####  Wallet

###### functions

`func GetEthKeySet(keystoreFilepath, password string) (ethKeySet EthereumKeyset, err error)`

`func (e *EthereumKeyset) SignTx(tx *types.Transaction, height string) (*types.Transaction, error)`

`func GenerateAndSaveEthereumWallet(keystoreFilepath, password string) error`

`func GenerateAndSaveEthereumWalletFromPrivateKey(privateKey, keystoreFilepath, password string) error`

`func (e *EthereumKeyset) ExportPrivateKey() (string, error)`


--------------------------------------------------------------------

### ABI

The ABI (application binary interface) along with the Archon Ethereum Smart Contract address are kept in the `abi/abi.go` file along with other minor utilities. For the most part, this file has no API's that would be useful to the developerand should be left alone. The ADC protocol relies on the Archon Ethereum Smart Contract as a control to various parts of the protocol. Said differently, each player in the ADC protocol acts in response to the evolving state of the Archon Ethereum Smart Contract in conjunction with other events. By design, it is not advantageous in the ADC protocol for an entity to manipulate its interface with the SC, or the SC itself (changing address in this file). The design effectively excises such entities from the protocol. "So don't bother".

--------------------------------------------------------------------

### Encodings

Like the ABI file above, the `encodings` folder contains assets that are necessary to this module, but are not intended to be used directly by the developer. Many of the API's exposed by this repository require various encodings/decodings behind the scenes. 

--------------------------------------------------------------------

