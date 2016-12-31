package blockchain

import (
	"reflect"

	"github.com/btcsuite/fastsha256"
	merkle "github.com/xsleonard/go-merkle"
)

// TxSlice contains a list of transactions
type TxSlice []*Tx

// Exists returns whether the tx exists in the slice
func (txs TxSlice) Exists(tx *Tx) bool {
	h := tx.Hash()
	for _, t := range txs {
		if reflect.DeepEqual(h, t.Hash()) {
			return true
		}
	}

	return false
}

// Last transaction in the slice.  It returns nil if there are 0 transactions.
func (txs TxSlice) Last() *Tx {
	if len(txs) > 0 {
		return txs[len(txs)-1]
	}
	return nil
}

// MerkleRoot hash of the transaction slice
func (txs TxSlice) MerkleRoot() ([]byte, error) {
	if len(txs) == 0 {
		return ZeroHash(), nil
	}

	// encode transactions.
	data := make([][]byte, len(txs))
	for i, tx := range txs {
		// same format used to calcuate the tx.Hash().  here data is provided as the merkel
		// lib takes in a hash function
		data[i] = concat(tx.DataHash(), tx.PrevHash, tx.Source, tx.Destination)
	}

	// actual merkel root calculation
	tree := merkle.NewTree()
	err := tree.Generate(data, fastsha256.New())
	if err == nil {
		return tree.Root().Hash, nil
	}
	return ZeroHash(), err
}

// Diff this slice with the given slice.  The algo assumes transactions are
// sorted (which may be too big of an assumption)
func (txs TxSlice) Diff(b TxSlice) (diff TxSlice) {
	lastj := 0
	for _, t := range txs {
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
