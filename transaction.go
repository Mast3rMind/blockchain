package blockchain

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/btcsuite/fastsha256"
	"github.com/xsleonard/go-merkle"
)

// Signator is used to sign a transaction
type Signator interface {
	Sign([]byte) (*Signature, []byte, error)
}

// TxType is a transaction type
type TxType uint8

// TxHeader contains header info for a transaction.
type TxHeader struct {
	prevHash   []byte
	pubKey     []byte
	nextPubKey []byte
}

// Tx represents a single transaction
type Tx struct {
	header    *TxHeader
	Signature *Signature
	data      []byte
}

// NewTx given the previous tx hash, data and optional public keys
func NewTx(prevHash []byte, data []byte) *Tx {
	return &Tx{
		header: &TxHeader{
			prevHash: prevHash,
		},
		data: data,
	}
}

// Header of the transaction
func (tx *Tx) Header() *TxHeader {
	return tx.header
}

// SetData for the tx
func (tx *Tx) SetData(d []byte) error {
	if tx.Signature != nil {
		return fmt.Errorf("tx already signed")
	}

	tx.data = d
	return nil
}

// Data of current tx
func (tx *Tx) Data() []byte {
	return tx.data
}

// DataHash of the tx data
func (tx *Tx) DataHash() []byte {
	s := fastsha256.Sum256(tx.data)
	return s[:]
}

// Hash of current data
func (tx *Tx) Hash() []byte {
	// data hash + previous hash + next pub key + signature (if signed)
	d := concat(tx.DataHash(), tx.header.prevHash, tx.header.nextPubKey)
	s := fastsha256.Sum256(d)
	return s[:]
}

// Sign transaction
func (tx *Tx) Sign(signer Signator) error {
	sig, pubkey, err := signer.Sign(tx.Hash())
	if err == nil {
		tx.Signature = sig
		tx.header.pubKey = pubkey
	}

	return err
}

func (tx *Tx) MarshalBinary() []byte {
	buf := new(bytes.Buffer)
	buf.Write(append(tx.Hash(), ' '))
	buf.Write(append(tx.header.prevHash, ' '))
	if tx.Signature != nil {
		buf.Write(tx.Signature.Bytes())
	}
	buf.Write([]byte{' '})
	buf.Write(tx.Data())
	return buf.Bytes()
}

func (tx *Tx) UnmarshalBinary(b []byte) error {
	buf := bytes.NewBuffer(b)
	hash, err := buf.ReadBytes(' ')
	if err != nil {
		return err
	}
	hash = hash[:len(hash)-1]

	if tx.header == nil {
		tx.header = &TxHeader{}
	}

	if tx.header.prevHash, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	tx.header.prevHash = tx.header.prevHash[:len(tx.header.prevHash)-1]

	var sig []byte
	if sig, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	if sig = sig[:len(sig)-1]; len(sig) > 0 {
		if tx.Signature, err = NewSignatureFromBytes(sig); err != nil {
			return err
		}
	}

	if tx.data, err = ioutil.ReadAll(buf); err != nil {
		return err
	}

	if !bytes.Equal(tx.Hash(), hash) {
		return fmt.Errorf("hash mistmatch")
	}
	return nil
}

// Verify the transaction signatures
func (tx *Tx) VerifySignature() (bool, error) {
	return tx.Signature.Verify(tx.header.pubKey, tx.Hash())
}

// TxSlice contains a list of transactions
type TxSlice []*Tx

func (txs TxSlice) Exists(tx *Tx) bool {
	h := tx.Hash()
	for _, t := range txs {
		if bytes.Equal(h, t.Hash()) {
			return true
		}
	}

	return false
}

// MerkleRoot hash of the transaction slice
func (txs TxSlice) MerkleRoot() ([]byte, error) {
	if len(txs) == 0 {
		return ZeroHash(), nil
	}

	// encode transactions.  May need to use DataHash instead of tx.data
	data := make([][]byte, len(txs))
	for i, tx := range txs {
		//data[i] = concat(tx.Hash(), []byte{' '}, tx.data)
		data[i] = tx.MarshalBinary()
	}

	// actual merkel root calculation
	tree := merkle.NewTree()
	if err := tree.Generate(data, fastsha256.New()); err != nil {
		return ZeroHash(), err
	}

	return tree.Root().Hash, nil
}
