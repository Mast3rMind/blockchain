package blockchain

import (
	"reflect"
	"testing"
	"time"
)

var (
	testBlkHdr1 = BlockHeader{
		Origin:     []byte("public-key"),
		PrevHash:   []byte("Previousblockhash"),
		MerkelRoot: []byte("merklerootofalltransactions"),
		Timestamp:  time.Now().UnixNano(),
		Nonce:      0,
	}
)

func Test_Block(t *testing.T) {
	txs := genSlices()
	blk := NewBlock(ZeroHash(), txs...)
	//h1 := blk.Header()
	m1 := blk.MerkelRoot

	if err := blk.AddTransaction(NewTx(blk.Transactions.Last().Hash(), []byte("dlkfajd;lkfjd;fkdjioeurpqiruewp"))); err != nil {
		t.Fatal(err)
	}

	if err := blk.AddTransaction(NewTx(ZeroHash(), []byte("dlkfajd;lkfjd;fkdjioeurpqiruewp"))); err == nil {
		t.Fatal("should fail with prev hash mismatch")
	}

	if reflect.DeepEqual(m1, blk.MerkelRoot) {
		t.Fatal("merkle roots should be different")
	}
}
