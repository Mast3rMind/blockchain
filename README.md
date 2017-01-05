# blockchain

blockchain provides an implementation of blockchain algorithm.  It is
meant to be a base for general purpose distributed transactions.

# QuickStart

To simply play around and get familiar with the library, a [CLI tool](https://github.com/ipkg/blockchain/tree/v2-dev/cmd)
is provided where transactions can be submitted via STDIN.

# Design

The implementation has the following 4 interfaces for customizations:

- StateMachine
- BlockStore
- Transport
- Signator

## StateMachine
This interface must be implemented by the user.  This is called each time a new
valid block is available.  It only contains one function that needs to be implemented.

A sample implementation can be found [here](https://github.com/ipkg/blockchain/tree/v2-dev/cmd).

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

- [ ] Re-work block request logic
- [ ] Add key management
- [ ] Optimize
- [ ] Increase test coverage


#### Ackknowledgments

The original blockchain code is based on [izqui/blockchain](https://github.com/izqui/blockchain) and
[gitchain/gitchain](https://github.com/gitchain/gitchain).
