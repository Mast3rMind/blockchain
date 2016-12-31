package blockchain

import (
	"bytes"
	"testing"
)

func Test_TxSlice(t *testing.T) {
	txs := genSlices()
	b, err := txs.MerkleRoot()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(ZeroHash(), b) {
		t.Fatal("should be non zero hash")
	}

	txs = TxSlice{}
	b, err = txs.MerkleRoot()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ZeroHash(), b) {
		t.Fatal("should be zero hash")
	}

	t.Logf("%x", b)
}
