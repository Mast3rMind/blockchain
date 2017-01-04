package blockchain

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/anacrolix/utp"
	chord "github.com/euforia/go-chord"
	"github.com/ipkg/go-mux"
)

var (
	testDialTimeout = time.Duration(1 * time.Second)
	testRpcTimeout  = time.Duration(50 * time.Millisecond)
	testMaxConnIdle = time.Duration(300 * time.Second)
)

type DummyTransport struct {
}

func (dt *DummyTransport) Initialize(tx chan<- *Tx, blk chan<- Block, store BlockStore) error {
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

func prepRingUTP(port int) (*mux.Mux, *chord.Config, *chord.UTPTransport, error) {
	listen := fmt.Sprintf("127.0.0.1:%d", port)
	conf := chord.DefaultConfig(listen)
	conf.StabilizeMin = time.Duration(15 * time.Millisecond)
	conf.StabilizeMax = time.Duration(45 * time.Millisecond)

	ln, err := utp.NewSocket("udp", listen)
	if err != nil {
		return nil, nil, nil, err
	}

	mx := mux.NewMux(ln, ln.Addr())
	go mx.Serve()
	sock1 := mx.Listen(72)

	trans, err := chord.InitUTPTransport(sock1, testDialTimeout, testRpcTimeout, testMaxConnIdle)
	if err != nil {
		return nil, nil, nil, err
	}
	return mx, conf, trans, nil
}

func Test_ChordTransport(t *testing.T) {
	// NODE 1

	mx1, cfg1, t1, err := prepRingUTP(40001)
	if err != nil {
		t.Fatal(err)
	}
	r1, err := chord.Create(cfg1, t1)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(200 * time.Millisecond)
	st1 := NewInMemBlockStore()
	bct1 := NewChordTransport(mx1.Listen(73), cfg1, r1, testDialTimeout, testMaxConnIdle)
	bc1, err := NewBlockchain(testKp, st1, bct1, &testFsm{})
	if err != nil {
		t.Fatal(err)
	}

	go bc1.Start()

	// NODE 2

	<-time.After(200 * time.Millisecond)
	mx2, cfg2, t2, err := prepRingUTP(40002)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := chord.Join(cfg2, t2, cfg1.Hostname)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(200 * time.Millisecond)
	// Set the genesis block from other store
	st2 := NewInMemBlockStore()
	st2.Add(*st1.LastBlock())
	bct2 := NewChordTransport(mx2.Listen(73), cfg2, r2, testDialTimeout, testMaxConnIdle)
	bc2, err := NewBlockchain(testKp, st2, bct2, &testFsm{}, cfg1.Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if st2.FirstBlock() == nil {
		t.Fatal("first block should not be nil")
	}

	go bc2.Start()

	<-time.After(200 * time.Millisecond)

	// CHECKS

	fb, err := bct2.FirstBlock(cfg1.Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if fb == nil {
		t.Fatal("first block should not be nil")
	}

	txs := genSlices()
	txs[0].Sign(testKp)
	txs[1].PrevHash = txs[0].Hash()
	txs[1].Sign(testKp)
	txs[2].PrevHash = txs[1].Hash()
	txs[2].Sign(testKp)

	bc1.QueueTransactions(txs...)

	<-time.After(3 * time.Second)

	b1 := bc1.store.LastBlock()
	if b1 == nil {
		t.Fatal("bc1 should have blocks")
	}

	b2 := bc2.store.LastBlock()
	if b2 == nil {
		t.Fatal("bc2 should have blocks")
	}

	if !reflect.DeepEqual(b1.Hash(), b2.Hash()) {
		t.Fatal("chain last block hash mismatch")
	}

	if bc1.store.BlockCount() != bc2.store.BlockCount() {
		t.Fatal("block count mismatch")
	}

	blb1, err := bct1.LastBlock(cfg1.Hostname)
	if err != nil {
		t.Fatal("bc1 failed to get last block")
	}
	blb2, err := bct2.LastBlock(cfg2.Hostname)
	if err != nil {
		t.Fatal("bc2 failed to get last block")
	}
	if !reflect.DeepEqual(blb1.Hash(), blb2.Hash()) {
		t.Fatal("last block mismatch between nodes")
	}

	// This should be the order of shutdown.
	r1.Leave() // leave ring
	<-time.After(1 * time.Second)
	//bc1.Shutdown()  // Shutdown blockchain
	t1.Shutdown()   // Shutdown ring transport
	bct1.Shutdown() // Shutdown chain transport

	<-time.After(1 * time.Second)
	r2.Leave() // leave ring
	<-time.After(1 * time.Second)
	//bc2.Shutdown()  // Shutdown blockchain
	t2.Shutdown()   // Shutdown ring transport
	bct2.Shutdown() // Shutdown chain transport

}
