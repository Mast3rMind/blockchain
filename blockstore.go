package blockchain

import (
	"log"
	"reflect"
)

type BlockStore interface {
	Exists(Block) bool
	Add(Block) error
	PreviousBlock() *Block
	Get(string) *Block
}

type InMemBlockStore struct {
	bs []Block
}

func NewInMemBlockStore() *InMemBlockStore {
	return &InMemBlockStore{bs: []Block{}}
}

func (ibs *InMemBlockStore) Exists(b Block) bool {
	//Traverse array in reverse order because if a block exists is more likely to be on top.
	l := len(ibs.bs)
	for i := l - 1; i >= 0; i-- {
		bb := ibs.bs[i]
		if reflect.DeepEqual(b.Signature, bb.Signature) {
			return true
		}
	}
	return false
}

func (ibs *InMemBlockStore) Get(hsh string) *Block {
	for _, b := range ibs.bs {
		if b.String() == hsh {
			return &b
		}
	}
	return nil
}

func (ibs *InMemBlockStore) Add(b Block) error {
	ibs.bs = append(ibs.bs, b)
	log.Printf("Chain size=%d", len(ibs.bs))
	return nil
}

func (ibs *InMemBlockStore) PreviousBlock() *Block {
	i := len(ibs.bs)
	if i == 0 {
		return nil
	}
	return &ibs.bs[i-1]
}
