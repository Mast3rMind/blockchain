package blockchain

import (
	"io"
	"log"
	"net"
	"time"
)

type NodeChannel chan *Node

type Nodes map[string]*Node

func (n Nodes) AddNode(node *Node) bool {
	key := node.TCPConn.RemoteAddr().String()
	//if key != Core.Network.Address && n[key] == nil {
	if _, ok := n[key]; !ok {
		log.Println("Node connected:", key)
		n[key] = node
		go node.handleNode()

		return true
	}
	return false
}

// Connection to node
type Node struct {
	*net.TCPConn
	lastSeen int
	// These are messages coming from the node.  They are read by the
	// local node.
	inMsg chan Message
}

func NewNode(conn *net.TCPConn, inChan chan Message) *Node {
	return &Node{
		TCPConn:  conn,
		inMsg:    inChan,
		lastSeen: int(time.Now().Unix()),
	}
}

func (node *Node) handleNode() {

	for {
		var bs []byte = make([]byte, 1024*1000)
		n, err := node.TCPConn.Read(bs[0:])
		if err != nil {
			if err == io.EOF {
				//TODO: Remove node [Issue: https://github.com/izqui/blockchain/issues/3]
				log.Println("Node disconnected:", node.RemoteAddr().String())
				if err = node.TCPConn.Close(); err != nil {
					log.Println("ERR", err)
				}
			} else {
				log.Println("ERR", err)
			}
			break
		}

		m := new(Message)
		if err = m.UnmarshalBinary(bs[0:n]); err != nil {
			log.Println("ERR", err)
			continue
		}

		m.Reply = make(chan Message)

		go func(cb chan Message) {
			for {
				m, ok := <-cb

				if !ok {
					close(cb)
					break
				}

				b, _ := m.MarshalBinary()
				l := len(b)

				i := 0
				for i < l {
					a, _ := node.TCPConn.Write(b[i:])
					i += a
				}
			}

		}(m.Reply)

		node.inMsg <- *m
		//Core.Network.IncomingMessages <- *m
	}
}
