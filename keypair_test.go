package blockchain

import (
	"crypto/sha256"
	"testing"

	"github.com/ipkg/blockchain/utils"
)

func TestKeyGeneration(t *testing.T) {
	keypair := GenerateNewKeypair()
	if len(keypair.Public) > 80 {
		t.Error("Error generating key")
	}
}

func TestKeySigning(t *testing.T) {

	for i := 0; i < 5; i++ {
		keypair := GenerateNewKeypair()

		data := utils.ArrayOfBytes(i, 'a')
		hash := sha256.Sum256(data)

		signature, err := keypair.Sign(hash[:])

		if err != nil {

			t.Error("base58 error")

		} else if !SignatureVerify(keypair.Public, signature, hash[:]) {

			t.Error("Signing and verifying error", len(keypair.Public))
		}
	}

}
