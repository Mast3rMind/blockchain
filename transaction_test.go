package blockchain

import (
	"bytes"
	"testing"
)

func genSlices() TxSlice {
	return TxSlice{
		NewTx(nil, []byte("werhgedfbrih6yvowtmeupcwipr")),
		NewTx(nil, []byte("weoq2cmpirotmrpj3imycwphlfkdsnjfgl;k")),
		NewTx(nil, []byte("etp7one56,buivmcorhoi3mj,j;")),
	}
}

func Test_Tx_Sign_Verify(t *testing.T) {
	tx := NewTx(ZeroHash(), nil)

	if err := tx.SetData([]byte("foobarbaz")); err != nil {
		t.Fatal(err)
	}
	kp, err := GenerateECDSAKeypair()
	if err != nil {
		t.Fatal(err)
	}

	if err = tx.Sign(kp); err != nil {
		t.Fatal(err)
	}

	if tx.Signature == nil {
		t.Fatal("failed to sign")
	}

	hdr := tx.Header()
	pkbytes := hdr.pubKey
	if pkbytes == nil || len(pkbytes) < 1 {
		t.Fatal("public key not set")
	}

	ok, err := tx.VerifySignature()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("failed to verify")
	}

	if err = tx.SetData(tx.Data()); err == nil {
		t.Fatal("should fail signature error")
	}

}

func Test_Tx_Marshal_Unmarshal(t *testing.T) {
	txs := genSlices()
	tx := txs[0]

	kpair, _ := GenerateECDSAKeypair()
	tx.Sign(kpair)

	b := tx.MarshalBinary()

	tx1 := &Tx{}
	if err := tx1.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(tx.Hash(), tx1.Hash()) {
		t.Fatal("hash mismatch")
	}

	if !bytes.Equal(tx.header.prevHash, tx1.header.prevHash) {
		t.Fatal("prev hash mismatch")
	}
	if !bytes.Equal(tx.Signature.Bytes(), tx1.Signature.Bytes()) {
		t.Fatal("signature mismatch")
	}

}

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
