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
	Timestamp  int64  `bencode:"t"`
	Nonce      uint32 `bencode:"n"`
	Bits       uint32 `bencode:"b"`
	Origin     []byte `bencode:"o"` // Public key of the node where the block originated from
}

// marshalBinary encodes the header to a bytes.  this is used for hash calculation
func (bh *BlockHeader) marshalBinary() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(append(bh.PrevHash, bh.MerkelRoot...))
	binary.Write(buf, binary.BigEndian, bh.Timestamp)
	binary.Write(buf, binary.BigEndian, bh.Nonce)
	binary.Write(buf, binary.BigEndian, bh.Bits)
	buf.Write(bh.Origin)
	return buf.Bytes()
}

/*// UnmarshalBinary decodes the bytes into the BlockHeader
func (bh *BlockHeader) UnmarshalBinary(b []byte) error {
	buf := bufio.NewReader(bytes.NewBuffer(b))
	return bh.decodeBuffer(buf)
}

// Decode data from reader into header
func (bh *BlockHeader) Decode(r io.Reader) error {
	rd := bufio.NewReader(r)
	return bh.decodeBuffer(rd)
}

func (bh *BlockHeader) decodeBuffer(rd *bufio.Reader) error {
	var err error

	if bh.PrevHash, err = rd.ReadBytes(' '); err != nil {
		return err
	}
	bh.PrevHash = bh.PrevHash[:len(bh.PrevHash)-1]

	if bh.MerkelRoot, err = rd.ReadBytes(' '); err != nil {
		return err
	}
	bh.MerkelRoot = bh.MerkelRoot[:len(bh.MerkelRoot)-1]

	if err = binary.Read(rd, binary.BigEndian, &bh.Timestamp); err == nil {
		if err = binary.Read(rd, binary.BigEndian, &bh.Nonce); err == nil {
			err = binary.Read(rd, binary.BigEndian, &bh.Bits)
		}
	}

	if err == nil {
		if bh.Origin, err = rd.ReadBytes(' '); err == nil {
			bh.Origin = bh.Origin[:len(bh.Origin)-1]
		} else if err == io.EOF {
			err = nil
		}
	}

	return err
}*/

// Hash of the encoded data
func (bh *BlockHeader) Hash() []byte {
	s := fastsha256.Sum256(bh.marshalBinary())
	return s[:]
}
