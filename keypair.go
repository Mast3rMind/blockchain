package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/tv42/base58"
)

const ecdsaKeySize = 28

// ECDSAKeypair is used to satisfy the keypair interface
type ECDSAKeypair struct {
	PrivateKey *ecdsa.PrivateKey
}

// GenerateECDSAKeypair generates a new keypair
func GenerateECDSAKeypair() (*ECDSAKeypair, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err == nil {
		return &ECDSAKeypair{PrivateKey: pk}, nil
	}
	return nil, err
}

// Sign data returning the 2 signatures
func (ekp *ECDSAKeypair) Sign(data []byte) (*Signature, []byte, error) {
	sr, ss, err := ecdsa.Sign(rand.Reader, ekp.PrivateKey, data)
	if err == nil {
		pub := ekp.PrivateKey.PublicKey
		pk := base58.EncodeBig([]byte{}, joinBigInt(ecdsaKeySize, pub.X, pub.Y))
		return &Signature{r: sr, s: ss}, pk, nil
		//encodeECDSAPublicKey(ekp.PrivateKey.PublicKey), nil
	}
	return nil, nil, err
}

func (ekp *ECDSAKeypair) PublicKeyBytes() []byte {
	x := ekp.PrivateKey.PublicKey.X
	y := ekp.PrivateKey.PublicKey.Y
	b := joinBigInt(ecdsaKeySize, x, y)
	return base58.EncodeBig([]byte{}, b)
}
