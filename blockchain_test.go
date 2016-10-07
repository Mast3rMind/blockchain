package blockchain

import (
	"reflect"
	"testing"
	//"time"
)

func Test_DiffTransactionSlices(t *testing.T) {
	bcTx1 := Transaction{Signature: randomData()}
	bcTx2 := Transaction{Signature: randomData()}
	bcTx3 := Transaction{Signature: randomData()}

	diff := DiffTransactionSlices(TransactionSlice{bcTx1, bcTx2, bcTx3}, TransactionSlice{bcTx1, bcTx3})
	if len(diff) != 1 && !reflect.DeepEqual(diff[0].Signature, bcTx2.Signature) {
		t.Error("Diffing algorithm fails")
	}
}

func randomData() []byte {
	return []byte(randomString(randomInt(0, 1024*1024)))
}

func buildTx(kp *Keypair) *Transaction {
	t := NewTransaction(kp.Public, nil, randomData())
	t.Header.Nonce = t.GenerateNonce(TRANSACTION_POW)
	t.Signature = t.Sign(kp)
	return t
}
