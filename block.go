package blockchain

import (
	"reflect"
	"time"
)

const (
	HIGHEST_TARGET = 0x1d00ffff
)

// Block consists of 1 or more transactions
type Block struct {
	*BlockHeader
	Signature    []byte  `bencode:"s"`
	Transactions TxSlice `bencode:"tx"`
}

// NewBlock instantiates a block with the previous hash and an initial set of
// optional transactions
func NewBlock(prevHash []byte, txs ...*Tx) *Block {
	b := &Block{
		BlockHeader: &BlockHeader{
			PrevHash:  prevHash,
			Timestamp: time.Now().UnixNano(),
			Nonce:     0,
			Bits:      HIGHEST_TARGET,
		},
		Transactions: TxSlice(txs),
	}
	if txs != nil {
		b.Transactions = txs
	}

	b.MerkelRoot, _ = b.Transactions.MerkleRoot()

	return b
}

// NewGenesisBlock returns a new genesis block i.e the previous hash is set to
// zero. One of these will exist per chain.
func NewGenesisBlock() *Block {
	return NewBlock(ZeroHash(), NewTx(ZeroHash(), []byte{}))
}

// Hash of the encoded header
func (blk *Block) Hash() []byte {
	return blk.BlockHeader.Hash()
}

// AddTransaction to the block
func (blk *Block) AddTransaction(tx *Tx) error {
	if blk.Signature != nil {
		return errAlreadySigned
	}

	ltx := blk.Transactions.Last()
	if ltx != nil {
		if !reflect.DeepEqual(ltx.Hash(), tx.PrevHash) {
			return errPrevHash
		}
	}

	txs := append(blk.Transactions, tx)
	mroot, err := txs.MerkleRoot()
	if err == nil {
		blk.Transactions = txs
		blk.MerkelRoot = mroot
	}

	return err
}

// Verify proof of work with the given prefix
func (blk *Block) Verify(prefix []byte) bool {
	headerHash := blk.Hash()
	merkel, _ := blk.Transactions.MerkleRoot()
	return reflect.DeepEqual(merkel, blk.MerkelRoot) && CheckProofOfWork(prefix, headerHash)
}

// VerifySignature of the block given the Signator to verify
func (blk *Block) VerifySignature(verifier Signator) error {
	return verifier.Verify(blk.Origin, blk.Signature, blk.Hash())
}

// Sign the block.
func (blk *Block) Sign(signer Signator) error {
	blk.Origin = signer.PublicKey().Bytes()

	sig, err := signer.Sign(blk.Hash())
	if err == nil {
		blk.Signature = sig.Bytes()
	}

	return err
}
