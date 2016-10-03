package blockchain

import (
	"log"
)

type core struct {
	*Keypair
	*Blockchain
	*Network

	cfg *Config
}

func NewCore(cfg *Config) *core {
	c := &core{cfg: cfg}
	c.setupKeypair()
	return c
}

func (c *core) setupKeypair() {
	keypair, _ := OpenConfiguration(c.cfg.DataDir)
	if keypair == nil {
		keypair = GenerateNewKeypair()
		log.Printf("Keypair generated. PublicKey: %x", keypair.Public)
		WriteConfiguration(c.cfg.DataDir, keypair)
	}
	c.Keypair = keypair
}

func (c *core) setupNetwork() {
	c.Network = NewNetwork(c.cfg.BindAddr)
	go c.Network.Run()
	// Add seeds
	for _, n := range c.cfg.SeedNodes {
		c.Network.ConnectionsQueue <- n
	}
}
func (c *core) setupBlockchain() {
	c.Blockchain = NewBlockchain(c.Keypair, c.Network.BroadcastQueue)
	go c.Blockchain.Run()
}

func (c *core) bootstrap() {
	c.setupNetwork()
	// must be called after network is setup
	c.setupBlockchain()
}

func (c *core) Start() {
	c.bootstrap()

	go func() {
		for {
			select {
			case msg := <-c.Network.IncomingMessages:
				c.handleIncomingMessage(msg)
			}
		}
	}()
}

func (c *core) handleIncomingMessage(msg Message) {

	switch msg.Identifier {
	case MESSAGE_SEND_TRANSACTION:
		t := new(Transaction)
		_, err := t.UnmarshalBinary(msg.Data)
		if err != nil {
			networkError(err)
			break
		}
		c.Blockchain.TransactionsQueue <- t

	case MESSAGE_SEND_BLOCK:
		b := new(Block)
		err := b.UnmarshalBinary(msg.Data)
		if err != nil {
			networkError(err)
			break
		}
		c.Blockchain.BlocksQueue <- *b

	default:
		log.Println("Unknown identifier:", msg.Identifier)

	}
}

func (c *core) CreateTransaction(payload []byte) *Transaction {
	t := NewTransaction(c.Keypair.Public, nil, payload)
	t.Header.Nonce = t.GenerateNonce(TRANSACTION_POW)
	t.Signature = t.Sign(c.Keypair)
	return t
}

// short to create and submit arbitrary data as transaction
func (c *core) SubmitTransaction(payload []byte) *Transaction {
	t := c.CreateTransaction(payload)
	c.Blockchain.TransactionsQueue <- t
	return t
}
