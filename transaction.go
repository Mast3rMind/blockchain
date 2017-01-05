package blockchain

import "github.com/btcsuite/fastsha256"

// Signator is used to sign a transaction
type Signator interface {
	Sign([]byte) (*Signature, error)
	PublicKey() PublicKey
	Verify(pubkey, signature, hash []byte) error
}

// TxType is a transaction type
type TxType uint8

// TxHeader contains header info for a transaction.
type TxHeader struct {
	PrevHash    []byte `bencode:"p"`
	Source      []byte `bencode:"f"` // from pubkey
	Destination []byte `bencode:"t"` // to pubkey
}

// Tx represents a single transaction
type Tx struct {
	*TxHeader
	//Signature *Signature `bencode:"s"`
	Signature []byte `bencode:"s"`
	Data      []byte `bencode:"d"`
}

// NewTx given the previous tx hash, data and optional public keys
func NewTx(prevHash []byte, data []byte) *Tx {
	return &Tx{
		TxHeader: &TxHeader{
			PrevHash: prevHash,
		},
		Data: data,
	}
}

// DataHash of the tx data
func (tx *Tx) DataHash() []byte {
	s := fastsha256.Sum256(tx.Data)
	return s[:]
}

// Hash of current data
func (tx *Tx) Hash() []byte {
	// data hash + previous hash + next pub key + signature (if signed)
	d := concat(tx.DataHash(), tx.PrevHash, tx.Source, tx.Destination)
	s := fastsha256.Sum256(d)
	return s[:]
}

// Sign transaction
func (tx *Tx) Sign(signer Signator) error {
	tx.Source = signer.PublicKey().Bytes()

	sig, err := signer.Sign(tx.Hash())
	if err == nil {
		tx.Signature = sig.Bytes()
		//tx.Source = pubkey
	}

	return err
}

// VerifySignature of the transaction
func (tx *Tx) VerifySignature(verifier Signator) error {
	return verifier.Verify(tx.Source, tx.Signature, tx.Hash())
}
