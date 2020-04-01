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

This module serves as an api interface for applications participating with the Archon Decentralized Cloud (ADC) to interact with the Archon Ethereum Smart Contract.


--------------------------------------------------------------------

### High-Level protocol description

For this simple protocol description, we define the players in the Archon Decentralized Cloud to be storage providers S, uploaders U, and downloaders D. The intent of these players is what you think it would be: the U want to make their conent available, the D want to obtain the content of U, and S wants to earn cryptocurrency by acting as a conduit serving the needs of U and D.

To bootstrap this protocol, we start with the S. Any storage provider S must be registerd with the Archon Ethereum Smart Contract as storage providers. This registrations includes providing information about their storage capabilities, routing information, as well as staking Ethereum token. An uploader U must also be registered with the SC. This registration includes establishing a namespace, and publishing the public key corresponding to their pseudo-identity.

We follow an upload u from U to its final target, the downloader D.


--------------------------------------------------------------------

### APIs 

--------------------------------------------------------------------

#### for the SP (storage provider)

--------------------------------------------------------------------

#### for the Uploader

--------------------------------------------------------------------

#### for the downloader

--------------------------------------------------------------------
