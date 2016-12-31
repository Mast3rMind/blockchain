# blockchain

blockchain provides a library implementation of the blockchain algorithm.  It is
meant to be a base for a general purpose distributed transaction log.

The implementation has the following 3 interfaces for customizations:

- Keypair
- BlockStore
- Transport

## Signatures
There are 2 types of signatures:

- Block
- Transaction

### Block Signature
Public key of the node where the block originated from.

### Transaction Signature
Public key of the entity issuing the transaction.
