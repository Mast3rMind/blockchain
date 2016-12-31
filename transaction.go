package blockchain

import "github.com/btcsuite/fastsha256"

// Signator is used to sign a transaction
type Signator interface {
	Sign([]byte) (*Signature, error)
	PublicKey() PublicKey
}

// TxType is a transaction type
type TxType uint8

// TxHeader contains header info for a transaction.
type TxHeader struct {
	PrevHash    []byte `bencode:"p"`
	Source      []byte `bencode:"f"` // from pubkey
	Destination []byte `bencode:"t"` // to pubkey
}

// Tx represents a single transaction
type Tx struct {
	*TxHeader
	//Signature *Signature `bencode:"s"`
	Signature []byte `bencode:"s"`
	Data      []byte `bencode:"d"`
}

// NewTx given the previous tx hash, data and optional public keys
func NewTx(prevHash []byte, data []byte) *Tx {
	return &Tx{
		TxHeader: &TxHeader{
			PrevHash: prevHash,
		},
		Data: data,
	}
}

// DataHash of the tx data
func (tx *Tx) DataHash() []byte {
	s := fastsha256.Sum256(tx.Data)
	return s[:]
}

// Hash of current data
func (tx *Tx) Hash() []byte {
	// data hash + previous hash + next pub key + signature (if signed)
	d := concat(tx.DataHash(), tx.PrevHash, tx.Source, tx.Destination)
	s := fastsha256.Sum256(d)
	return s[:]
}

// Sign transaction
func (tx *Tx) Sign(signer Signator) error {
	tx.Source = signer.PublicKey().Bytes()

	sig, err := signer.Sign(tx.Hash())
	if err == nil {
		tx.Signature = sig.Bytes()
		//tx.Source = pubkey
	}

	return err
}

/*func (tx *Tx) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := tx.encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (tx *Tx) encode(w io.Writer) error {

	w.Write(append(tx.Hash(), ' '))
	w.Write(append(tx.PrevHash, ' '))

	if tx.Signature != nil {
		w.Write(tx.Signature)
	}
	w.Write([]byte{' '})

	w.Write(append(tx.Source, ' '))
	w.Write(append(tx.Destination, ' '))

	binary.Write(w, binary.LittleEndian, uint64(len(tx.Data)))
	w.Write(tx.Data)
	return nil
}


// UnmarshalBinary data into tx
func (tx *Tx) UnmarshalBinary(b []byte) error {
	buf := bufio.NewReader(bytes.NewBuffer(b))
	return tx.decoderBuffer(buf)
}

func (tx *Tx) Encode(w io.Writer) error {
	return tx.encode(w)
}

// Decode tx from reader
func (tx *Tx) Decode(r io.Reader) error {
	buf := bufio.NewReader(r)
	return tx.decoderBuffer(buf)
}

func (tx *Tx) decoderBuffer(buf *bufio.Reader) error {
	hash, err := buf.ReadBytes(' ')
	if err != nil {
		return err
	}
	hash = hash[:len(hash)-1]

	if tx.TxHeader == nil {
		tx.TxHeader = &TxHeader{}
	}

	if tx.PrevHash, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	tx.PrevHash = tx.PrevHash[:len(tx.PrevHash)-1]
	log.Printf("PREV %s", tx.PrevHash)

	var sig []byte
	if sig, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	if sig = sig[:len(sig)-1]; len(sig) > 0 {
		s1, err := NewSignatureFromBytes(sig)
		if err != nil {
			return err
		}
		tx.Signature = s1.Bytes()
	}

	if tx.Source, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	tx.Source = tx.Source[:len(tx.Source)-1]

	if tx.Destination, err = buf.ReadBytes(' '); err != nil {
		return err
	}
	tx.Destination = tx.Destination[:len(tx.Destination)-1]

	log.Printf("TX %#v", tx.TxHeader)

	sz, err := binary.ReadUvarint(buf)
	if err != nil {
		return err
	}
	//log.Println("SIZE", sz)

	tx.Data = make([]byte, sz)
	rn, err := io.ReadFull(buf, tx.Data)
	if err != nil {
		return err
	}
	if uint64(rn) != sz {
		return fmt.Errorf("data not complete")
	}

	if !bytes.Equal(tx.Hash(), hash) {
		log.Printf("FIX %#v", *tx.TxHeader)

		return fmt.Errorf("hash mistmatch")
	}

	return nil
}*/

// VerifySignature of the transaction
func (tx *Tx) VerifySignature() error {
	sig, err := NewSignatureFromBytes(tx.Signature)
	if err == nil {
		return sig.Verify(tx.Source, tx.Hash())
	}
	return err
}
