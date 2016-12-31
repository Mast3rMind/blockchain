# blockchain

blockchain provides an implementation of blockchain algorithm.  It is
meant to be a base for general purpose distributed transactions.

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
