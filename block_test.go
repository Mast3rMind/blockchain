package blockchain

import (
	"bytes"
	"testing"
	"time"
)

var (
	testBlkHdr1 = BlockHeader{
		PreviousHash: []byte("Previousblockhash"),
		MerkelRoot:   []byte("merklerootofalltransactions"),
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	}
)

func Test_BlockHeader(t *testing.T) {
	bh := testBlkHdr1
	ed := bh.Encode()

	bh1 := BlockHeader{}
	if err := bh1.Decode(ed); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bh.PreviousHash, bh1.PreviousHash) {
		t.Fatalf("prev hash mismatch '%x'!='%x'", bh.PreviousHash, bh1.PreviousHash)
	}

	if !bytes.Equal(bh.MerkelRoot, bh1.MerkelRoot) {
		t.Fatal("merkle root hash mismatch")
	}
}

func Test_Block(t *testing.T) {
	txs := genSlices()
	blk := NewBlock(ZeroHash(), txs)
	h1 := blk.Header()
	m1 := h1.MerkelRoot

	if err := blk.AddTransaction(NewTx(ZeroHash(), []byte("dlkfajd;lkfjd;fkdjioeurpqiruewp"))); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(m1, blk.header.MerkelRoot) {
		t.Fatal("merkle roots should be different")
	}
}
