package blockchain

import (
	"log"
	"reflect"
	"time"
)

// Transaction poll interval
const (
	defaultTxPollSecs int = 30
)

type Blockchain struct {
	CurrentBlock Block
	// Holds all available blocks
	BlockStore
	// Channel to read incoming transactions made available
	// to the internal engine for processing, verification etc.
	tq chan *Tx
	// Transaction exposore outbound
	// Transactions verfied from tq to be braodcasted.
	txOut chan *Tx
	// Channel on which generated blocks will be available
	// to the internal engine for processing verifcation etc.
	bq chan Block
	// Block exposure outbound
	// Verified blocks to be broadcasted
	blkOut chan Block

	// Public-Private keypair for this chain
	Keypair *ECDSAKeypair

	// Check for new transactions per the below interval
	// This is a constant check even it interrupt not received
	TxPollInterval int
}

//func NewBlockchain(keypair *Keypair, broadcastChan chan rpc.Message) *Blockchain {
func NewBlockchain(keypair *ECDSAKeypair, store BlockStore) *Blockchain {
	bl := &Blockchain{
		tq:             make(chan *Tx),
		bq:             make(chan Block),
		txOut:          make(chan *Tx),
		blkOut:         make(chan Block),
		Keypair:        keypair,
		TxPollInterval: defaultTxPollSecs,
	}
	if store == nil {
		bl.BlockStore = NewInMemBlockStore()
	} else {
		bl.BlockStore = store
	}

	// TODO: Read blockchain from file and stuff...
	bl.CurrentBlock = bl.createNewBlock()

	return bl
}

// QueueTransaction queues transaction to be added to a block.  Transactions are
// verified and then only are acccepted.
func (bl *Blockchain) QueueTransaction(tx *Tx) {
	bl.tq <- tx
}

// QueueBlock queues a block to be added to the chain.  This performs a full
// validation before actually adding it to the chain.
func (bl *Blockchain) QueueBlock(b Block) {
	bl.bq <- b
}

// BlockAvailable returns a channel containing valid and verified blocks that are
// available for broadcast on the network.
func (bl *Blockchain) BlockAvailable() <-chan Block {
	return bl.blkOut
}

// TransactionAvailable returns a channel containing valid and verified transactions
// that are available for broadcast on the network.
func (bl *Blockchain) TransactionAvailable() <-chan *Tx {
	return bl.txOut
}

func (bl *Blockchain) createNewBlock() Block {

	prevBlock := bl.PreviousBlock()
	prevBlockHash := []byte{}
	if prevBlock != nil {
		prevBlockHash = prevBlock.Hash()
	}

	b := NewBlock(prevBlockHash, nil)
	b.header.Origin = bl.Keypair.PublicKeyBytes()

	return *b
}

func (bl *Blockchain) Run() {

	interruptBlockGen := bl.generateBlocks()
	for {
		select {
		// Process transaction
		case tr := <-bl.tq:
			if bl.CurrentBlock.Transactions.Exists(tr) {
				continue
			}
			// verify signature only
			ok, err := tr.VerifySignature()
			if err != nil {
				log.Println(err)
				continue
			}
			if !ok {
				log.Printf("Verfication failed tx=%x", tr.Hash())
			}

			bl.CurrentBlock.AddTransaction(tr)
			interruptBlockGen <- bl.CurrentBlock
			// make transaction available to broadcast or do whatever else
			bl.txOut <- tr

		case b := <-bl.bq:
			// Process block
			if bl.Exists(b) {
				continue
			}
			ok, err := b.VerifySignature()
			if err != nil {
				log.Printf("Signature verfication failed: block=%x reason='%s'", b.Hash(), err.Error())
				continue
			}
			if !ok {
				log.Printf("Signature invalid: block=%x", b.Hash())
				continue
			}

			if !b.Verify(BLOCK_POW) {
				log.Printf("Verification failed: block=%x", b.Hash())
				continue
			}

			if reflect.DeepEqual(b.header.PreviousHash, bl.CurrentBlock.Hash()) {
				// I'm missing some blocks in the middle. Request'em.
				log.Printf("Missing blocks between prev=%x curr=%x", b.header.PreviousHash, bl.CurrentBlock.Hash())

			} else {
				transDiff := TxSlice{}
				if !reflect.DeepEqual(b.header.MerkelRoot, bl.CurrentBlock.header.MerkelRoot) {
					log.Println("Transactions are different. Calculating diff")
					transDiff = DiffTransactionSlices(bl.CurrentBlock.Transactions, b.Transactions)
				}

				log.Printf("Adding block=%x tx=%d", b.Hash(), len(b.Transactions))
				if e := bl.Add(b); e != nil {
					log.Println("ERR", e)
				}

				// Make block available to broadcast or do whatever else
				bl.blkOut <- b

				// Reset current block to a new block
				bl.CurrentBlock = bl.createNewBlock()
				bl.CurrentBlock.Transactions = transDiff
				interruptBlockGen <- bl.CurrentBlock
			}
		}
	}
}

// Start generating blocks
func (bl *Blockchain) generateBlocks() chan Block {
	interrupt := make(chan Block)

	go func() {
		block := <-interrupt

	loop:
		log.Printf("[POW] Begin block=%x", block.Hash())
		block.header.MerkelRoot, _ = block.Transactions.MerkleRoot()
		block.header.Nonce = 0
		block.header.Timestamp = time.Now().UnixNano()

		for {

			sleepTime := time.Nanosecond
			if len(block.Transactions) > 0 {
				//log.Println("[generateBlocks] Transactions", block.TransactionSlice.Len(), block.Nonce)
				if CheckProofOfWork(BLOCK_POW, block.Hash()) {
					log.Printf("[POW] Found block=%x", block.Hash())
					if err := block.Sign(bl.Keypair); err != nil {
						log.Println("Failed to sign block:", err)
					} else {
						//log.Printf("Signed block=%x tx=%d", block.Hash(), len(block.Transactions))
						bl.bq <- block
						// TODO: change to milliseconds
						sleepTime = time.Second * time.Duration(bl.TxPollInterval)
					}
				} else {
					block.header.Nonce++
				}

			} else {
				// TODO: change to milliseconds
				sleepTime = time.Second * time.Duration(bl.TxPollInterval)
				//log.Println("DBG [POW] No transactions. Sleeping for", sleepTime.Seconds(), "secs")
			}

			select {
			case block = <-interrupt:
				goto loop

			case <-time.After(sleepTime):
				continue
			}
		}
	}()

	return interrupt
}

// DiffTransactionSlices - Assumes transaction arrays are sorted (which may be too big of an assumption)
func DiffTransactionSlices(a, b TxSlice) (diff TxSlice) {
	lastj := 0
	for _, t := range a {
		found := false
		for j := lastj; j < len(b); j++ {
			if reflect.DeepEqual(b[j].Signature, t.Signature) {
				found = true
				lastj = j
				break
			}
		}

		if !found {
			diff = append(diff, t)
		}
	}

	return
}
