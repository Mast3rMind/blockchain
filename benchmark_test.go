package blockchain

import (
	"testing"
)

func BenchmarkTransaction(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doTransaction()
	}

}
func BenchmarkBlock(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doBlock()
	}
}

func doTransaction() {
	t := NewTransaction(nil, nil, nil)
	t.GenerateNonce(TRANSACTION_POW)
	t.Sign(GenerateNewKeypair())
}

func doBlock() {
	b := NewBlock(nil)
	b.GenerateMerkelRoot()
	b.GenerateNonce(BLOCK_POW)
	b.Sign(GenerateNewKeypair())
}
