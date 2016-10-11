package blockchain

import (
	"log"
	"reflect"
	"time"
)

// Transaction poll interval
const (
	DEFAULT_TX_POLL_TIME_SECS int = 30
	NETWORK_KEY_SIZE              = 80
)

type Blockchain struct {
	CurrentBlock Block
	// Holds all available blocks
	BlockStore
	// Channel to read incoming transactions made available
	// to the internal engine for processing, verification etc.
	tq chan *Transaction
	// Channel on which generated blocks will be available
	// to the internal engine for processing verifcation etc.
	bq chan Block
	// Transaction exposore outbound
	// Transactions verfied from tq to be braodcasted.
	txOut chan *Transaction
	// Block exposure outbound
	// Verified blocks to be broadcasted
	blkOut chan Block

	// Public-Private keypair for this chain
	Keypair *Keypair
	// Check for new transactions per the below interval
	// This is a constant check even it interrupt not received
	TxPollInterval int
}

//func NewBlockchain(keypair *Keypair, broadcastChan chan rpc.Message) *Blockchain {
func NewBlockchain(keypair *Keypair, store BlockStore) *Blockchain {
	bl := &Blockchain{
		tq:             make(chan *Transaction),
		bq:             make(chan Block),
		txOut:          make(chan *Transaction),
		blkOut:         make(chan Block),
		Keypair:        keypair,
		TxPollInterval: DEFAULT_TX_POLL_TIME_SECS,
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

// Queue transaction to be added to a block.  Transactions are verified and
// then only are acccepted.
func (bl *Blockchain) QueueTransaction(tx *Transaction) {
	bl.tq <- tx
}

// Queue a block to be added to the chain.  This performs a full
// validation before actually adding it to the chain.
func (bl *Blockchain) QueueBlock(b Block) {
	bl.bq <- b
}

// Valid and verified blocks that are available
func (bl *Blockchain) BlockAvailable() <-chan Block {
	return bl.blkOut
}

// Valid and verified transactions that are available
func (bl *Blockchain) TransactionAvailable() <-chan *Transaction {
	return bl.txOut
}

func (bl *Blockchain) createNewBlock() Block {

	//prevBlock := bl.BlockSlice.PreviousBlock()
	prevBlock := bl.PreviousBlock()
	prevBlockHash := []byte{}
	if prevBlock != nil {
		prevBlockHash = prevBlock.Hash()
	}

	b := NewBlock(prevBlockHash)
	b.BlockHeader.Origin = bl.Keypair.Public

	return b
}

func (bl *Blockchain) Run() {

	interruptBlockGen := bl.generateBlocks()
	for {
		select {
		case tr := <-bl.tq:
			if bl.CurrentBlock.TransactionSlice.Exists(*tr) {
				continue
			}
			if !tr.VerifyTransaction(TRANSACTION_POW) {
				log.Printf("Verfication failed tx=%s", tr.String())
				continue
			}

			bl.CurrentBlock.AddTransaction(tr)
			interruptBlockGen <- bl.CurrentBlock
			// make transaction available to broadcast or do whatever else
			bl.txOut <- tr

		case b := <-bl.bq:
			//if bl.BlockSlice.Exists(b) {
			if bl.Exists(b) {
				//log.Println("Exists block=%s", b.String())
				continue
			}
			if !b.VerifyBlock(BLOCK_POW) {
				log.Printf("Verification failed block=%s", b.String())
				continue
			}

			if reflect.DeepEqual(b.PrevBlock, bl.CurrentBlock.Hash()) {
				// I'm missing some blocks in the middle. Request'em.
				log.Printf("Missing blocks between prev=%x curr=%s", b.PrevBlock, bl.CurrentBlock.String())
			} else {
				transDiff := TransactionSlice{}
				if !reflect.DeepEqual(b.BlockHeader.MerkelRoot, bl.CurrentBlock.MerkelRoot) {
					log.Println("Transactions are different. Calculating diff")
					transDiff = DiffTransactionSlices(*bl.CurrentBlock.TransactionSlice, *b.TransactionSlice)
				}

				log.Printf("Adding block=%s", b.String())
				//bl.BlockSlice = append(bl.BlockSlice, b)
				if e := bl.Add(b); e != nil {
					log.Println("ERR", e)
				}

				// make block available to broadcast or do whatever else
				bl.blkOut <- b

				// Reset current block to a new block
				bl.CurrentBlock = bl.createNewBlock()
				bl.CurrentBlock.TransactionSlice = &transDiff

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
		log.Printf("[POW] Begin block=%s", block.String())
		block.BlockHeader.MerkelRoot = block.GenerateMerkelRoot()
		block.BlockHeader.Nonce = 0
		block.BlockHeader.Timestamp = uint32(time.Now().Unix())

		for {

			sleepTime := time.Nanosecond
			if block.TransactionSlice.Len() > 0 {
				//log.Println("[generateBlocks] Transactions", block.TransactionSlice.Len(), block.Nonce)
				if CheckProofOfWork(BLOCK_POW, block.Hash()) {
					log.Printf("[POW] Found block=%s", block.String())

					block.Signature = block.Sign(bl.Keypair)
					bl.bq <- block

					sleepTime = time.Second * time.Duration(bl.TxPollInterval)

				} else {
					block.BlockHeader.Nonce += 1
				}

			} else {
				//sleepTime = time.Hour * 24
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

//Assumes transaction arrays are sorted (which maybe is too big of an assumption)
func DiffTransactionSlices(a, b TransactionSlice) (diff TransactionSlice) {
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
