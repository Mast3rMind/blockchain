package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ipkg/blockchain"
)

func init() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	kp := blockchain.GenerateNewKeypair()
	chain := blockchain.NewBlockchain(kp, nil)
	go chain.Run()

	go func() {
		for {
			select {
			case tx := <-chain.TransactionAvailable():
				log.Printf("tx=%x", tx.Hash())

			case blk := <-chain.BlockAvailable():
				log.Printf("block=%x", blk.Hash())

			}
		}
	}()

	fmt.Println("Type something and hit enter")
	for {
		str := <-readStdin()

		tx := buildTx(kp, []byte(str))
		chain.QueueTransaction(tx)

		log.Printf("Submitted tx=%s", tx.String())
	}
}

func buildTx(kp *blockchain.Keypair, payload []byte) *blockchain.Transaction {
	tx := blockchain.NewTransaction(kp.Public, nil, payload)
	tx.Header.Nonce = tx.GenerateNonce(blockchain.TRANSACTION_POW)
	tx.Signature = tx.Sign(kp)
	return tx
}

func readStdin() chan string {

	cb := make(chan string)
	sc := bufio.NewScanner(os.Stdin)

	go func() {
		if sc.Scan() {
			cb <- sc.Text()
		}
	}()

	return cb
}
