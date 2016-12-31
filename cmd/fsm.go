package main

import (
	"log"

	"github.com/ipkg/blockchain"
)

type stateMachine struct {
}

func (sm *stateMachine) Apply(b blockchain.Block) error {
	log.Printf("[fsm.Apply] %x blk=%x nonce=%d tx=%d\n", b.PrevHash[len(b.PrevHash)-4:],
		b.Hash(), b.Nonce, len(b.Transactions))

	return nil
}
