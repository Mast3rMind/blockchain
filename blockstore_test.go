package blockchain

import (
	"bytes"
	"testing"
)

func Test_BlockStore(t *testing.T) {
	bs := NewInMemBlockStore()
	if bs.LastBlock() != nil {
		t.Fatal("last block should be nil")
	}

	if bs.FirstBlock() != nil {
		t.Fatal("first block should be nil")
	}

	if bs.LastTx() != nil {
		t.Fatal("last tx should be nil")
	}

	blk := NewGenesisBlock()
	bs.Add(*blk)

	if !bytes.Equal(bs.FirstBlock().Hash(), blk.Hash()) {
		t.Fatal("wrong first block")
	}

	lb := bs.LastBlock()
	if !bytes.Equal(lb.Hash(), blk.Hash()) {
		t.Fatal("wrong last block")
	}

	if bs.Get(blk.Hash()) == nil {
		t.Error("block should exist")
	}

	if bs.Get([]byte("ldkfjd")) != nil {
		t.Error("block should not exist")
	}

	tx1 := blk.Transactions[len(blk.Transactions)-1]
	tx2 := bs.LastTx()

	if !bytes.Equal(tx1.Hash(), tx2.Hash()) {
		t.Error("tx hash mismatch")
	}
}
