package blockchain

import (
	"log"
	"reflect"
	"testing"
	"time"
)

var (
	testKp, _ = GenerateECDSAKeypair()
)

type testFsm struct{}

func (tf *testFsm) Apply(b Block) error {
	log.Printf("fsm.Apply> blk=%x nonce=%d tx=%d\n", b.Hash(), b.Nonce, len(b.Transactions))
	return nil
}

func Test_DiffTransactionSlices(t *testing.T) {
	txs1 := genSlices()
	txs1[0].Sign(testKp)
	txs1[1].Sign(testKp)
	txs1[2].Sign(testKp)

	diff := txs1.Diff(txs1[:len(txs1)-1])
	if len(diff) != 1 &&
		!reflect.DeepEqual(diff[0].Signature, txs1[len(txs1)-1].Signature) {

		t.Error("Diffing algorithm fails")
	}
}

func Test_Blockchain(t *testing.T) {
	trans := &DummyTransport{}
	chain, err := NewBlockchain(testKp, nil, trans, &testFsm{})
	if err != nil {
		t.Fatal(err)
	}
	go chain.Start()

	txs := genSlices()

	txs[0].Sign(testKp)
	txs[1].PrevHash = txs[0].Hash()
	txs[1].Sign(testKp)
	txs[2].PrevHash = txs[1].Hash()
	txs[2].Sign(testKp)

	chain.QueueTransactions(txs...)

	<-time.After(2 * time.Second)
}
