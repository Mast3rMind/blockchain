package blockchain

import (
	"bytes"
	"sync"
)

// InMemBlockStore is an in memory blockstore
type InMemBlockStore struct {
	mu sync.Mutex
	bs []Block
}

// NewInMemBlockStore initialized with the genisis block with zero hash
func NewInMemBlockStore() *InMemBlockStore {
	return &InMemBlockStore{bs: []Block{}}
}

// Exists returns if the store has the given block. It traverses the array in
// reverse order as likelyhood of newwer blocks is on top.
func (ibs *InMemBlockStore) Exists(b Block) bool {

	l := len(ibs.bs)
	for i := l - 1; i >= 0; i-- {
		bb := ibs.bs[i]
		if bytes.Equal(b.Hash(), bb.Hash()) && bytes.Equal(b.PrevHash, bb.PrevHash) {
			return true
		}
	}
	return false
}

// Get block from the store given its hash
func (ibs *InMemBlockStore) Get(hsh []byte) *Block {
	for _, b := range ibs.bs {
		if bytes.Equal(b.Hash(), hsh) {
			return &b
		}
	}
	return nil
}

// Add a block to the store it if doesn't exist
func (ibs *InMemBlockStore) Add(blk Block) error {
	ibs.mu.Lock()
	defer ibs.mu.Unlock()

	ibs.bs = append(ibs.bs, blk)
	//log.Printf("Chain size=%d", len(ibs.bs))
	return nil
}

// NewBlock with the previous hash set to the hash of LastBlock
func (ibs *InMemBlockStore) NewBlock() *Block {
	prevBlock := ibs.LastBlock()
	return NewBlock(prevBlock.Hash())
}

// LastBlock in the store
func (ibs *InMemBlockStore) LastBlock() *Block {
	i := len(ibs.bs)
	if i == 0 {
		return nil
	}
	return &ibs.bs[i-1]
}

// LastTx returns the last tx in the last block.
func (ibs *InMemBlockStore) LastTx() *Tx {
	blk := ibs.LastBlock()
	if blk == nil {
		return nil
	}

	if blk.Transactions == nil || len(blk.Transactions) < 1 {
		return nil
	}

	return blk.Transactions.Last()
}

// BlockCount for the store
func (ibs *InMemBlockStore) BlockCount() int64 {
	return int64(len(ibs.bs))
}

// FirstBlock returns the first block in this store. ie. chain
func (ibs *InMemBlockStore) FirstBlock() *Block {
	return &ibs.bs[0]
}
