package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"github.com/tv42/base58"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func NewSignatureFromBytes(b []byte) (*Signature, error) {
	bi, err := base58.DecodeToBig(b)
	if err != nil {
		return nil, err
	}

	pp := splitBigInt(bi, 2)
	return &Signature{r: pp[0], s: pp[1]}, nil
}

// Bytes of encoded signature
func (sig *Signature) Bytes() []byte {
	b := joinBigInt(ecdsaKeySize, sig.r, sig.s)
	return base58.EncodeBig([]byte{}, b)
}

// Verify data given the public key using the signature
func (sig *Signature) Verify(pubkey []byte, data []byte) (bool, error) {
	pub, err := decodeECDSAPublicKeyBytes(pubkey)
	if err == nil {
		return ecdsa.Verify(&pub, data, sig.r, sig.s), nil
	}
	return false, err
}

func decodeECDSAPublicKeyBytes(pk []byte) (pub ecdsa.PublicKey, err error) {
	var b *big.Int
	if b, err = base58.DecodeToBig(pk); err == nil {
		sigg := splitBigInt(b, 2)
		pub = ecdsa.PublicKey{Curve: elliptic.P256(), X: sigg[0], Y: sigg[1]}
	}

	return
}
