package main

import (
	"log"

	"github.com/ipkg/blockchain"
)

type stateMachine struct{}

// Apply the block.  This is where all user work should be done.  Returning an
// error causes the block to be rejected.
func (sm *stateMachine) Apply(b blockchain.Block) error {
	log.Printf("[fsm.Apply] %x blk=%x nonce=%d tx=%d\n", b.PrevHash[len(b.PrevHash)-4:],
		b.Hash(), b.Nonce, len(b.Transactions))

	return nil
}
