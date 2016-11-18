# blockchain
This is a blockchain only library assembled and modified from [izqui/blockchain](https://github.com/izqui/blockchain).  The networking component is not included.

### Implementation
It is meant to be a standalone library and hence the network implementation has been removed.  Instead channels are used to feed and extract transactions and blocks to and from the chain.

This allows for custom implementaions of the networking layer such as using Chord, Pastry, Kademlia, or even your own custom RPC mechanism.

### Example

	kp := blockchain.GenerateNewKeypair()
	chain := blockchain.NewBlockchain(kp)
	go chain.Run()

	// Start reading blocks and transactions that have been verified, validated and added
	// to the chain. 
	go func() {
		for {
			select {
			// Channel with available transactions from the chain
			case tx := <-chain.TransactionAvailable():
				log.Printf("tx=%x", tx.Hash())

			// Channel with available blocks from the chain
			case bk := <-chain.BlockAvailable():
				log.Printf("block=%x", bk.Hash())

			}
		}
	}()

	// Generate a transaction every 3 seconds
	for {
		<-time.After(3)
		
		tx := blockchain.NewTransaction(kp.Public, nil, []byte("some-data"))
		tx.Header.Nonce = tx.GenerateNonce(blockchain.TRANSACTION_POW)
		tx.Signature = tx.Sign(kp)
		// Queue the transaction.  As the transaction gets validated and a block is
		// generated, it be available on the output channels
		chain.QueueTransaction(tx)
	}

### Notes

Critic and improvements are always welcome!