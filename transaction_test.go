package blockchain

import "testing"

func genSlices() TxSlice {
	return TxSlice{
		NewTx(nil, []byte("werhgedfbrih6yvowtmeupcwipr")),
		NewTx(nil, []byte("weoq2cmpirotmrpj3imycwphlfkdsnjfgl;k")),
		NewTx(nil, []byte("etp7one56,buivmcorhoi3mj,j;")),
	}
}

func Test_Tx_Sign_Verify(t *testing.T) {
	tx := NewTx(ZeroHash(), nil)

	tx.Data = []byte("foobarbaz")

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

	pkbytes := tx.Source
	if pkbytes == nil || len(pkbytes) < 1 {
		t.Fatal("public key not set")
	}

	if err = tx.VerifySignature(kp); err != nil {
		t.Fatal(err)
	}

}
