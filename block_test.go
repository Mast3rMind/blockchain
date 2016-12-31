package blockchain

import (
	"bytes"
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

/*func Test_BlockHeader(t *testing.T) {
	bh := testBlkHdr1
	ed := bh.MarshalBinary()

	bh1 := &BlockHeader{}
	if err := bh1.UnmarshalBinary(ed); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bh.Hash(), bh1.Hash()) {
		t.Fatalf("hash mismatch '%x'!='%x'", bh.Hash(), bh1.Hash())
	}

	if !bytes.Equal(bh.PrevHash, bh1.PrevHash) {
		t.Fatalf("prev hash mismatch '%x'!='%x'", bh.PrevHash, bh1.PrevHash)
	}

	if !bytes.Equal(bh.MerkelRoot, bh1.MerkelRoot) {
		t.Fatal("merkle root hash mismatch")
	}

	if !bytes.Equal(bh.Origin, bh1.Origin) {
		t.Fatal("merkle root hash mismatch")
	}
}*/

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

	if bytes.Equal(m1, blk.MerkelRoot) {
		t.Fatal("merkle roots should be different")
	}
}
