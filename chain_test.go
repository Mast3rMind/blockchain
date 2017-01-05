package blockchain

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

var (
	testKp, _ = GenerateECDSAKeypair()
)

type DummyTransport struct {
}

func (dt *DummyTransport) Initialize(tx chan<- *Tx, blk chan<- Block, store ReadOnlyBlockStore) error {
	return nil
}

func (dt *DummyTransport) BroadcastTransaction(*Tx) error {
	return fmt.Errorf("tbi")
}

func (dt *DummyTransport) BroadcastBlock(*Block) error {
	return fmt.Errorf("tbi")
}

func (dt *DummyTransport) RequestBlocks(hashes ...[]byte) {
}

func (dt *DummyTransport) FirstBlock(host string) (*Block, error) {
	return nil, fmt.Errorf("tbi")
}
func (dt *DummyTransport) LastBlock(host string) (*Block, error) {
	return nil, fmt.Errorf("tbi")
}

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
