package main

import (
	"fmt"
	"time"

	"github.com/ipkg/blockchain"
)

func main() {

	fmt.Println("Benchmarking...")

	t1 := benchmark(func() {
		t := blockchain.NewTransaction(nil, nil, nil)
		t.GenerateNonce(blockchain.TRANSACTION_POW)
		t.Sign(blockchain.GenerateNewKeypair())
	})
	fmt.Println("Transaction took", t1)

	t2 := benchmark(func() {
		b := blockchain.NewBlock(nil)
		b.GenerateMerkelRoot()
		b.GenerateNonce(blockchain.BLOCK_POW)
		b.Sign(blockchain.GenerateNewKeypair())
	})
	fmt.Println("Block took", t2)
}

func benchmark(f func()) time.Duration {

	t0 := time.Now()

	f()

	return time.Since(t0)
}
