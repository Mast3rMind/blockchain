package blockchain

import (
	"bytes"
	"errors"

	"github.com/ipkg/blockchain/utils"
)

const (
	MESSAGE_TYPE_SIZE    = 1
	MESSAGE_OPTIONS_SIZE = 4
)

const (
	MESSAGE_GET_NODES = iota + 20
	MESSAGE_SEND_NODES

	MESSAGE_GET_TRANSACTION
	MESSAGE_SEND_TRANSACTION

	MESSAGE_GET_BLOCK
	MESSAGE_SEND_BLOCK
)

type Message struct {
	Identifier byte
	Options    []byte
	Data       []byte

	Reply chan Message
}

func NewMessage(id byte) *Message {
	return &Message{Identifier: id}
}

func (m *Message) MarshalBinary() ([]byte, error) {

	buf := new(bytes.Buffer)

	buf.WriteByte(m.Identifier)
	buf.Write(utils.FitBytesInto(m.Options, MESSAGE_OPTIONS_SIZE))
	buf.Write(m.Data)

	return buf.Bytes(), nil

}

func (m *Message) UnmarshalBinary(d []byte) error {
	if len(d) < MESSAGE_OPTIONS_SIZE+MESSAGE_TYPE_SIZE {
		return errors.New("insuficient message size")
	}

	buf := bytes.NewBuffer(d)

	m.Identifier = buf.Next(1)[0]
	m.Options = utils.StripByte(buf.Next(MESSAGE_OPTIONS_SIZE), 0)
	m.Data = buf.Next(utils.MaxInt)

	return nil
}
