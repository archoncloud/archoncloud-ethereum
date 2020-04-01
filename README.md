# archoncloud-ethereum

NOTE: This software is in development and subject to change. Not for production use!

Documentation IN PROGRESS.

### Contents:

 - Overview

 - High-Level protocol description

 - APIs 

   - for the SP (storage provider)

   - for The Uploader

   - for The Downloader

--------------------------------------------------------------------


### Overview

This module serves as an api interface for applications interacting with the Archon Decentralized Cloud (ADC) either as agents of the cloud, or as beneficiaries. Specifically, this module is an api interface to interact with the Archon Ethereum Smart Contract.


--------------------------------------------------------------------

### High-Level protocol description

This is a very high-level description. Many details are glossed over in order to keep this brief. For a more detailed description, see the Archon White Paper /TODO LINK/ or read the source code from our official repositories.

For this simple protocol description, we define the players in the Archon Decentralized Cloud to be storage providers S, uploaders U, and downloaders D. The intent of these players are what you think they would be: the U want to make their content available, the D want to obtain the content of U, and S wants to earn cryptocurrency by acting as a conduit serving the needs of U and D.

To bootstrap this protocol, we start with the S. Any storage provider S must be registered with the Archon Ethereum Smart Contract as storage providers. This registrations includes providing information about their storage capabilities, marketplace ask, routing information, as well as staking Ethereum token. An uploader U must also be registered with the SC. This registration includes establishing a namespace, and publishing the public key corresponding to their pseudo-identity.

We follow an upload u from U to its final target, the downloader D.

To keep this description very simple, we abstract away the details of erasure-encoding, etc and just follow u from U to S to D.

The U prepares u using some encoding and signs u to get {u,sig(u), <other-metadata>}. Either now, or in the past, U has accumulated a subset of storage providers S_ = {S_1,S_2,...,S_n} (a subset of all storage providers in ADC) from one or a few S. Locally, U runs the ADC marketplace to determine the best S from S_ to accomodate u.

U concurrently makes a proposeUpload transaction pu_tx to SC with metadata of u and S that was determined by the marketplace and sends {u, sig(u), <other-metadata>}. S caches this upload and listens to SC for pu_tx to be confirmed by the blockchain. The proposeUpload transaction includes a payment to S for storing u, documents metadata of u including sig(u), and also validates the result of the marketplace. Assuming pu_tx is confirmed, S announces to the networking overlay of ADC that it is storing {u, ...} and stores {u,...} for the period paid for by U in pu_tx. 

The downloader D knows of u from some other channel. Perhaps U advertised on, say, reddit that U uploaded u. D contacts some storage provider S' asking for the ADC download url of u. Storage provider S' queries its networking overlay for the url(s) of any S holding u and returns these values to D. Downloader D downloads {u, sig(u), <other-metadata>} from S and retrieves the public key of U from the SC. D validates sig(u) with this public key and accepts u in the ideal case.

We will see below which API's each of the players call in order to participate in this protocol. Please keep in mind, this description glossed over some very important implementation details in order to be brief. For a more detailed protocol description, refer to the Archon Whitepaper /TODO NEED URL/, or inspect the source of the official repositories.

--------------------------------------------------------------------

### APIs 

--------------------------------------------------------------------

#### for the SP (storage provider)

--------------------------------------------------------------------

#### for the Uploader

--------------------------------------------------------------------

#### for the downloader

--------------------------------------------------------------------
