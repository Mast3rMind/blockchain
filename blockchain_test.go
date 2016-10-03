package blockchain

import (
	"reflect"
	"testing"
)

func TestBlockDiff(t *testing.T) {

	tr1 := Transaction{Signature: []byte(randomString(randomInt(0, 1024*1024)))}
	tr2 := Transaction{Signature: []byte(randomString(randomInt(0, 1024*1024)))}
	tr3 := Transaction{Signature: []byte(randomString(randomInt(0, 1024*1024)))}

	diff := DiffTransactionSlices(TransactionSlice{tr1, tr2, tr3}, TransactionSlice{tr1, tr3})

	if len(diff) != 1 && !reflect.DeepEqual(diff[0].Signature, tr2.Signature) {

		t.Error("Diffing algorithm fails")
	}
}
