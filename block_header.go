package blockchain

import (
	"bytes"
	"encoding/binary"

	"github.com/btcsuite/fastsha256"
)

// BlockHeader contains block metadata
type BlockHeader struct {
	PrevHash   []byte `bencode:"p"` // Previous block hash
	MerkelRoot []byte `bencode:"m"` // Merkel root hash of all transactions
	Timestamp  int64  `bencode:"t"` // Block creation time
	Nonce      uint32 `bencode:"n"`
	Bits       uint32 `bencode:"b"`
	Origin     []byte `bencode:"o"` // Public key of the node where the block originated from
}

// marshalBinary encodes the header to bytes.  This is used for hash calculation
// of the header
func (bh *BlockHeader) marshalBinary() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(append(bh.PrevHash, bh.MerkelRoot...))
	binary.Write(buf, binary.BigEndian, bh.Timestamp)
	binary.Write(buf, binary.BigEndian, bh.Nonce)
	binary.Write(buf, binary.BigEndian, bh.Bits)
	buf.Write(bh.Origin)
	return buf.Bytes()
}

// Hash of the encoded data
func (bh *BlockHeader) Hash() []byte {
	s := fastsha256.Sum256(bh.marshalBinary())
	return s[:]
}
