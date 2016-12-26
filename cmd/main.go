package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ipkg/blockchain"
)

func buildTx(kp *blockchain.ECDSAKeypair, prevHash, payload []byte) (*blockchain.Tx, error) {
	//tx := blockchain.NewTransaction(kp.Public, nil, payload)
	tx := blockchain.NewTx(prevHash, payload)
	//tx.Header.Nonce = tx.GenerateNonce(stx.TRANSACTION_POW)
	err := tx.Sign(kp)
	return tx, err
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

func main() {
	kp, err := blockchain.GenerateECDSAKeypair()
	if err != nil {
		log.Fatal(err)
	}
	store := blockchain.NewInMemBlockStore()

	chain := blockchain.NewBlockchain(kp, store)
	go chain.Run()

	prevHash := blockchain.ZeroHash()

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

		tx, err := buildTx(kp, prevHash, []byte(str))
		if err != nil {
			log.Println("ERR", err)
			continue
		}
		chain.QueueTransaction(tx)

		log.Printf("Submitted tx=%x", tx.Hash())
		prevHash = tx.Hash()
	}
}
