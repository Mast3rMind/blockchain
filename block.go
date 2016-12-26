package blockchain

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/btcsuite/fastsha256"
)

const (
	HIGHEST_TARGET = 0x1d00ffff
)

// BlockHeader contains block metadata
type BlockHeader struct {
	Origin       []byte // public key of node
	PreviousHash []byte // Previous block hash
	MerkelRoot   []byte // merkle root of all transactions
	Timestamp    int64
	Nonce        uint32
	Bits         uint32
}

// Encode the block header.  This does not contain the signature
func (bh *BlockHeader) Encode() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(append(bh.Origin, ' '))
	buf.Write(append(bh.PreviousHash, ' '))
	buf.Write(append(bh.MerkelRoot, ' '))

	binary.Write(buf, binary.BigEndian, bh.Timestamp)
	binary.Write(buf, binary.BigEndian, bh.Nonce)
	binary.Write(buf, binary.BigEndian, bh.Bits)

	return buf.Bytes()
}

// Decode the bytes into the BlockHeader
func (bh *BlockHeader) Decode(b []byte) error {
	var (
		rd  = bufio.NewReader(bytes.NewBuffer(b))
		err error
	)

	if bh.Origin, err = rd.ReadBytes(' '); err != nil {
		return err
	}
	bh.Origin = bh.Origin[:len(bh.Origin)-1]

	if bh.PreviousHash, err = rd.ReadBytes(' '); err != nil {
		return err
	}
	bh.PreviousHash = bh.PreviousHash[:len(bh.PreviousHash)-1]

	if bh.MerkelRoot, err = rd.ReadBytes(' '); err != nil {
		return err
	}
	bh.MerkelRoot = bh.MerkelRoot[:len(bh.MerkelRoot)-1]

	if err = binary.Read(rd, binary.BigEndian, bh.Timestamp); err == nil {
		if err = binary.Read(rd, binary.BigEndian, bh.Nonce); err == nil {
			err = binary.Read(rd, binary.BigEndian, bh.Bits)
		}
	}
	return err
}

// Hash of the encoded data
func (bh *BlockHeader) Hash() []byte {
	s := fastsha256.Sum256(bh.Encode())
	return s[:]
}

// Block consists of 1 or more transactions
type Block struct {
	//mu     sync.Mutex
	Signature    *Signature
	header       *BlockHeader
	Transactions TxSlice
}

// NewBlock instantiates a block with the previous hash and an initial set of
// optional transactions
func NewBlock(prevHash []byte, txs TxSlice) *Block {
	b := &Block{
		header: &BlockHeader{
			PreviousHash: prevHash,
			Timestamp:    time.Now().UnixNano(),
			Nonce:        0,
			Bits:         HIGHEST_TARGET,
		},
		Transactions: TxSlice{},
	}
	if txs != nil {
		b.Transactions = txs
	}

	b.header.MerkelRoot, _ = b.Transactions.MerkleRoot()

	return b
}

// Header of the block reod-only
func (blk *Block) Header() BlockHeader {
	return *blk.header
}

// Hash of the encoded header
func (blk *Block) Hash() []byte {
	return blk.header.Hash()
}

// AddTransaction to the block
func (blk *Block) AddTransaction(tx *Tx) error {
	if blk.Signature != nil {
		return fmt.Errorf("block already signed")
	}
	//blk.Transactions = append(blk.Transactions, tx)
	txs := append(blk.Transactions, tx)
	mroot, err := txs.MerkleRoot()
	if err == nil {
		//blk.mu.Lock()
		blk.Transactions = txs
		blk.header.MerkelRoot = mroot
		//blk.mu.Unlock()
	}

	return err
}

func (blk *Block) Verify(prefix []byte) bool {
	headerHash := blk.Hash()
	merkel, _ := blk.Transactions.MerkleRoot()
	return reflect.DeepEqual(merkel, blk.header.MerkelRoot) &&
		CheckProofOfWork(prefix, headerHash)
	//SignatureVerify(b.BlockHeader.Origin, b.Signature, headerHash)
}

func (blk *Block) VerifySignature() (bool, error) {
	if blk.Signature == nil {
		return false, fmt.Errorf("block not signed: %x", blk.Hash())
	}
	return blk.Signature.Verify(blk.header.Origin, blk.Hash())
}

func (blk *Block) Sign(signer Signator) error {
	sig, pubkey, err := signer.Sign(blk.Hash())
	if err == nil {
		blk.Signature = sig
		blk.header.Origin = pubkey
	}

	return err
}
