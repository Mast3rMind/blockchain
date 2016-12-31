# blockchain

blockchain provides an implementation of blockchain algorithm.  It is
meant to be a base for general purpose distributed transactions.

The implementation has the following 3 interfaces for customizations:

- BlockStore
- Transport
- Signator

## BlockStore
A an in-memory block store is provided called **InMemBlockStore**

## Transport
An implementation of a chord based transport is provided that uses uTP called **ChordTransport**

## Signator
The signator interface is responsible for signing and verifying blocks and transactions.  The provided
implementation uses ECDSA key pairs called **ECDSAKeypair**.

### Block Signature
This should be the public key of the node where the block originated from.

### Transaction Signature
This should be the public key of the entity issuing the transaction.

#### To Do

- [ ] Key management
- [ ] Increase test coverage


#### Ackknowledgments

The original blockchain code is based on [izqui/blockchain](https://github.com/izqui/blockchain) and
[gitchain/gitchain](https://github.com/gitchain/gitchain).
