package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"
)

const TRANSACTION_HEADER_SIZE = NETWORK_KEY_SIZE /* from key */ + NETWORK_KEY_SIZE /* to key */ + 4 /* int32 timestamp */ + 32 /* sha256 payload hash */ + 4 /* int32 payload length */ + 4 /* int32 nonce */

type Transaction struct {
	Header    TransactionHeader
	Signature []byte
	Payload   []byte
}

type TransactionHeader struct {
	From          []byte
	To            []byte
	Timestamp     uint32
	PayloadHash   []byte
	PayloadLength uint32
	Nonce         uint32
}

// Returns bytes to be sent to the network
func NewTransaction(from, to, payload []byte) *Transaction {
	t := &Transaction{
		Header: TransactionHeader{
			From:      from,
			To:        to,
			Timestamp: uint32(time.Now().Unix()),
		},
		Payload: payload,
	}

	sh := sha256.Sum256(t.Payload)
	t.Header.PayloadHash = sh[:]

	t.Header.PayloadLength = uint32(len(t.Payload))
	return t
}

func (t *Transaction) Hash() []byte {
	headerBytes, _ := t.Header.MarshalBinary()
	sh := sha256.Sum256(headerBytes)
	return sh[:]
}

// Returns string representation of hash
func (t *Transaction) String() string {
	return fmt.Sprintf("%x", t.Hash())
}

func (t *Transaction) Sign(keypair *Keypair) []byte {
	s, _ := keypair.Sign(t.Hash())
	return s
}

func (t *Transaction) VerifyTransaction(pow []byte) bool {
	headerHash := t.Hash()
	sh := sha256.Sum256(t.Payload)
	payloadHash := sh[:]

	return reflect.DeepEqual(payloadHash, t.Header.PayloadHash) &&
		CheckProofOfWork(pow, headerHash) &&
		SignatureVerify(t.Header.From, t.Signature, headerHash)
}

func (t *Transaction) GenerateNonce(prefix []byte) uint32 {
	newT := t
	for {
		if CheckProofOfWork(prefix, newT.Hash()) {
			break
		}
		newT.Header.Nonce++
	}
	return newT.Header.Nonce
}

func (t *Transaction) MarshalBinary() ([]byte, error) {
	headerBytes, _ := t.Header.MarshalBinary()

	if len(headerBytes) != TRANSACTION_HEADER_SIZE {
		return nil, errors.New("Header marshalling error")
	}

	return append(append(headerBytes, fitBytesInto(t.Signature, NETWORK_KEY_SIZE)...), t.Payload...), nil
}

func (t *Transaction) UnmarshalBinary(d []byte) ([]byte, error) {
	buf := bytes.NewBuffer(d)

	if len(d) < TRANSACTION_HEADER_SIZE+NETWORK_KEY_SIZE {
		return nil, errors.New("Insuficient bytes for unmarshalling transaction")
	}

	header := &TransactionHeader{}
	if err := header.UnmarshalBinary(buf.Next(TRANSACTION_HEADER_SIZE)); err != nil {
		return nil, err
	}
	t.Header = *header

	t.Signature = stripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	t.Payload = buf.Next(int(t.Header.PayloadLength))

	return buf.Next(MAX_DATA_SIZE), nil
}

func (th *TransactionHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(fitBytesInto(th.From, NETWORK_KEY_SIZE))
	buf.Write(fitBytesInto(th.To, NETWORK_KEY_SIZE))
	binary.Write(buf, binary.LittleEndian, th.Timestamp)
	buf.Write(fitBytesInto(th.PayloadHash, 32))
	binary.Write(buf, binary.LittleEndian, th.PayloadLength)
	binary.Write(buf, binary.LittleEndian, th.Nonce)

	return buf.Bytes(), nil

}

func (th *TransactionHeader) UnmarshalBinary(d []byte) error {
	buf := bytes.NewBuffer(d)

	th.From = stripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	th.To = stripByte(buf.Next(NETWORK_KEY_SIZE), 0)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.Timestamp)
	th.PayloadHash = buf.Next(32)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.PayloadLength)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &th.Nonce)

	return nil
}

type TransactionSlice []Transaction

func (slice TransactionSlice) Len() int {
	return len(slice)
}

func (slice TransactionSlice) Exists(tr Transaction) bool {
	for _, t := range slice {
		if reflect.DeepEqual(t.Signature, tr.Signature) {
			return true
		}
	}
	return false
}

func (slice TransactionSlice) AddTransaction(t Transaction) TransactionSlice {
	// Inserted sorted by timestamp
	for i, tr := range slice {
		if tr.Header.Timestamp >= t.Header.Timestamp {
			return append(append(slice[:i], t), slice[i:]...)
		}
	}

	return append(slice, t)
}

func (slice *TransactionSlice) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, t := range *slice {
		bs, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(bs)
	}

	return buf.Bytes(), nil
}

func (slice *TransactionSlice) UnmarshalBinary(d []byte) error {
	remaining := d
	for len(remaining) > TRANSACTION_HEADER_SIZE+NETWORK_KEY_SIZE {
		t := new(Transaction)
		rem, err := t.UnmarshalBinary(remaining)

		if err != nil {
			return err
		}
		(*slice) = append((*slice), *t)
		remaining = rem
	}
	return nil
}
