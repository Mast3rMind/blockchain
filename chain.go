package blockchain

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"
)

// Transaction poll interval
const (
	defaultTxPollSecs int = 10
)

var (
	errPrevHash = fmt.Errorf("previous hash mismatch")
)

// Transport protocol for the network.  This abstracts out the peers.
type Transport interface {
	// write only channel for network blocks and txs to be queued to the local chain
	Initialize(chan<- *Tx, chan<- Block, BlockStore) error
	BroadcastTransaction(*Tx) error
	BroadcastBlock(*Block) error
	// RequestBlocks will requests the blocks with the given hashes and write them
	// to the block channel as provided in Initialize.  This will be called from as
	// a go routine.
	RequestBlocks(hashes ...[]byte)
	// Last block in the chain as seen by hosts local state
	LastBlock(host string) (*Block, error)
	// First block in the chain as seen by this node.  This may not be the genesis
	// block as it will be different foreach node depending on when they joined
	// the tranasction chain.
	FirstBlock(host string) (*Block, error)
}

// BlockStore contains all the valid mined blocks that have been processed for a
// single chain
type BlockStore interface {
	// new block with the previous hash.
	NewBlock() *Block
	// check if block exists
	Exists(Block) bool
	// add block to store once validated
	Add(Block) error
	// first block in the chain
	FirstBlock() *Block
	// last block in the chain
	LastBlock() *Block
	// last transaction in the chain
	LastTx() *Tx
	// get block by hash
	Get(hash []byte) *Block
	// total block count in the store
	BlockCount() int64
}

// StateMachine is the user implemented state machine.  Apply is called each time
// a new block is available i.e. mined, verified and approved to be added to the
// chain.
type StateMachine interface {
	// Apply is called each time a new valid block is made available.  Returning
	// an error will cause the block not be added to the chain.
	Apply(Block) error
}

type Blockchain struct {
	// lock for the current
	mu sync.Mutex
	// block being worked on currently
	curBlk Block

	// Channel to read incoming transactions made available
	// to the internal engine for processing, verification etc.
	tq chan *Tx
	// Channel on which generated blocks will be available
	// to the internal engine for processing verifcation etc.
	bq chan Block
	// Channel to generate blocks
	genBlk chan Block

	// Holds all available blocks
	store BlockStore

	// implemented interface to handle mined blocks
	fsm StateMachine

	// Public-Private keypair for this chain.  The public key is used as Origin
	// for every block generated
	signator Signator

	// netowrk transport
	transport Transport

	// Check for new transactions per the below interval
	// This is a constant check even it interrupt not received
	TxPollInterval int
	// Time to wait for before moving onto the next tx loop
	//TxTimeoutMsec int
}

// NewBlockchain initializes a new blockchain.  If the store does not contain an
// existing last block, a new genesis block is created.  The transport is initialized
// with the tx and block channels of the chain so txs and blocks from the network
// can be received.
func NewBlockchain(signator Signator, store BlockStore, trans Transport, fsm StateMachine, peers ...string) (*Blockchain, error) {
	bl := &Blockchain{
		tq:             make(chan *Tx),
		bq:             make(chan Block),
		genBlk:         make(chan Block),
		signator:       signator,
		TxPollInterval: defaultTxPollSecs,
		store:          store,
		transport:      trans,
		fsm:            fsm,
	}

	if bl.store == nil {
		bl.store = NewInMemBlockStore()
	}

	// Initial supplied transport with tx and block channels
	// TODO: pass store to Initialize
	err := bl.transport.Initialize(bl.tq, bl.bq, store)
	if err != nil {
		return nil, err
	}

	if err = bl.initBlockStore(peers...); err != nil {
		return nil, err
	}

	bl.curBlk = bl.createNewBlock()
	return bl, nil
}

// QueueTransactions queues transaction to be added to a block.  Transactions are
// verified and then only are acccepted.
func (bl *Blockchain) QueueTransactions(tx ...*Tx) {
	for _, t := range tx {
		bl.tq <- t
	}
}

// createNewBlock for generation
func (bl *Blockchain) createNewBlock() Block {

	// TODO: may need to get last block from network.
	b := bl.store.NewBlock()
	b.Origin = bl.signator.PublicKey().Bytes()

	return *b
}

// verify transaction, and broadcast on success
func (bl *Blockchain) processTx(tx *Tx) error {
	if bl.curBlk.Transactions.Exists(tx) {
		return nil
	}
	// verify signature only
	err := tx.VerifySignature()
	if err != nil {
		return err
	}

	// TODO verify other parts of the tx

	// make transaction available to broadcast or do whatever else.  this should
	// be done first before we add to our current block
	if e := bl.transport.BroadcastTransaction(tx); e != nil {
		log.Printf("ERR [transport] broadcast tx=%x reason='%v'", tx.Hash(), e)
	}

	// add tx to curr block
	bl.mu.Lock()
	if err = bl.curBlk.AddTransaction(tx); err != nil {
		bl.mu.Unlock()
		return err
	}
	bl.mu.Unlock()

	// send curr block to block generation channel
	bl.genBlk <- bl.curBlk

	return nil
}

// this runs in a go routine and may need to be optimized as many without limit
// can be spawned
func (bl *Blockchain) submitBlocksRequest(hashes ...[]byte) {
	arr := [][]byte{}
	for _, h := range hashes {
		if h != nil && len(h) > 0 {
			arr = append(arr, h)
		}
	}
	bl.transport.RequestBlocks(arr...)
}

// verify block, perform user callback, store block, broadcast to network
func (bl *Blockchain) processBlock(b Block) error {
	if bl.store.Exists(b) {
		return nil
	}

	err := b.VerifySignature()
	if err != nil {
		return err
	}

	// 3. Verfiy proof of work
	if !b.Verify(BLOCK_POW) {
		return fmt.Errorf("proof-of-work verification failed")
	}

	// We are missing blocks between b.PrevHash and bl.curBlk.Hash().  Request them
	// from the network.
	if bytes.Equal(bl.curBlk.Hash(), b.PrevHash) {
		if bl.store.Get(b.PrevHash) != nil {
			log.Printf("Chain may have diverged at: %x", b.PrevHash)
			return nil
		}

		log.Printf("Requesting block: %x", b.PrevHash)
		go bl.submitBlocksRequest(b.PrevHash)

		return nil
	}

	txDiff := TxSlice{}
	if !bytes.Equal(b.MerkelRoot, bl.curBlk.MerkelRoot) {
		txDiff = bl.curBlk.Transactions.Diff(b.Transactions)
		log.Printf("Transaction diff: %d", len(txDiff))
	}

	// Call user implemented StateMachine ie. any user specified work.  Once that
	// succeeds we add the block to the store otherwise the block is discarded.
	if err = bl.fsm.Apply(b); err != nil {
		return err
	}

	if err = bl.store.Add(b); err != nil {
		return err
	}

	// Make block available to broadcast or do whatever else
	if e := bl.transport.BroadcastBlock(&b); e != nil {
		log.Printf("ERR [transport] broadcast block=%x reason='%v'", b.Hash(), e)
	}

	// Reset current block to a new block
	bl.mu.Lock()
	bl.curBlk = bl.createNewBlock()
	bl.curBlk.Transactions = txDiff
	bl.mu.Unlock()

	bl.genBlk <- bl.curBlk

	return nil
}

// Start generating blocks
func (bl *Blockchain) startBlockGeneration() {
	block := <-bl.genBlk

loop:

	block.MerkelRoot, _ = block.Transactions.MerkleRoot()
	block.Nonce = 0
	block.Timestamp = time.Now().UnixNano()

	for {
		sleepTime := time.Nanosecond
		if len(block.Transactions) < 1 {
			sleepTime = time.Second * time.Duration(bl.TxPollInterval)
			goto enditer
		}

		if !CheckProofOfWork(BLOCK_POW, block.Hash()) {
			block.Nonce++
			goto enditer
		}

		if err := block.Sign(bl.signator); err != nil {
			log.Println("Failed to sign block:", err)
			goto enditer
		}

		bl.bq <- block

		// TODO: change to milliseconds
		sleepTime = time.Second * time.Duration(bl.TxPollInterval)

	enditer:

		select {
		case block = <-bl.genBlk:
			goto loop

		case <-time.After(sleepTime):
			continue
		}

	}
}

func (bl *Blockchain) Start() {
	go bl.startBlockGeneration()

	for {
		select {

		case tr := <-bl.tq:
			// Process transaction
			bl.processTx(tr)
			//if err := bl.processTx(tr); err != nil {
			//	log.Printf("ERR tx=%x reason='%s'", tr.Hash(), err.Error())
			//}

		case b := <-bl.bq:
			// Process block
			bl.processBlock(b)
			//if err := bl.processBlock(b); err != nil {
			//	log.Printf("ERR block=%x reason='%s'", b.Hash(), err.Error())
			//}

		}
	}
}

// If creating a new cluster check if blockstore has last block otherwise create a new
// genesis block as the last block.  If joining a cluster, ask the existing node for the
// last block.
//
// Bootstrap is an existing peer list.  If none are provided it assumes
// you want to create a new ring.
func (bl *Blockchain) initBlockStore(bootstrap ...string) error {
	// store has a last block. nothing to do. this is the case whether we're creating
	// or joining
	prevBlk := bl.store.LastBlock()
	if prevBlk != nil {
		return nil
	}

	var err error
	if len(bootstrap) == 0 {
		// Create cluster - Add genesis block if we're creating ring.
		prevBlk = NewGenesisBlock()
		prevBlk.Origin = bl.signator.PublicKey().Bytes()
	} else {
		// Joining a cluster - Get last block from existing node in cluster.  Try
		// all peers until a success.
		for _, p := range bootstrap {
			if prevBlk, err = bl.transport.LastBlock(p); err != nil {
				continue
			}
			break
		}
	}
	if err == nil {
		return bl.store.Add(*prevBlk)
	}
	return err
}
