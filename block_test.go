package blockchain

import (
	"crypto/rand"
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestMerkellHash(t *testing.T) {

	tr1 := NewTransaction(nil, nil, []byte(randomString(randomInt(0, 1024*1024))))
	tr2 := NewTransaction(nil, nil, []byte(randomString(randomInt(0, 1024*1024))))
	tr3 := NewTransaction(nil, nil, []byte(randomString(randomInt(0, 1024*1024))))
	tr4 := NewTransaction(nil, nil, []byte(randomString(randomInt(0, 1024*1024))))

	b := new(Block)
	b.TransactionSlice = &TransactionSlice{*tr1, *tr2, *tr3, *tr4}

	mt := b.GenerateMerkelRoot()
	manual := SHA256(append(SHA256(append(tr1.Hash(), tr2.Hash()...)), SHA256(append(tr3.Hash(), tr4.Hash()...))...))

	if !reflect.DeepEqual(mt, manual) {
		t.Error("Merkel tree generation fails")
	}
}

//TODO: Write block validation and marshalling tests [Issue: https://github.com/izqui/blockchain/issues/2]
func TestBlockMarshalling(t *testing.T) {

	kp := GenerateNewKeypair()
	tr := NewTransaction(kp.Public, nil, []byte(randomString(randomInt(0, 1024*1024))))

	tr.Header.Nonce = tr.GenerateNonce(arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX))
	tr.Signature = tr.Sign(kp)

	data, err := tr.MarshalBinary()

	if err != nil {
		t.Error(err)
	}

	newT := &Transaction{}
	_, err = newT.UnmarshalBinary(data)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*newT, *tr) {
		t.Error("Marshall, unmarshall failed")
	}
}

func TestBlockVerification(t *testing.T) {

	pow := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX)

	kp := GenerateNewKeypair()
	tr := NewTransaction(kp.Public, nil, []byte(randomString(randomInt(0, 1024))))

	tr.Header.Nonce = tr.GenerateNonce(pow)
	tr.Signature = tr.Sign(kp)

	if !tr.VerifyTransaction(pow) {

		t.Error("Validation failing")
	}
}

func TestIncorrectBlockPOWVerification(t *testing.T) {

	pow := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX)
	powIncorrect := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, 'a')

	kp := GenerateNewKeypair()
	tr := NewTransaction(kp.Public, nil, []byte(randomString(randomInt(0, 1024))))
	tr.Header.Nonce = tr.GenerateNonce(powIncorrect)
	tr.Signature = tr.Sign(kp)

	if tr.VerifyTransaction(pow) {

		t.Error("Passed validation without pow")
	}
}

func TestIncorrectBlockSignatureVerification(t *testing.T) {

	pow := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX)
	kp1, kp2 := GenerateNewKeypair(), GenerateNewKeypair()
	tr := NewTransaction(kp2.Public, nil, []byte(randomString(randomInt(0, 1024))))
	tr.Header.Nonce = tr.GenerateNonce(pow)
	tr.Signature = tr.Sign(kp1)

	if tr.VerifyTransaction(pow) {

		t.Error("Passed validation with incorrect key")
	}
}

// From http://devpy.wordpress.com/2013/10/24/create-random-string-in-golang/
func randomString(n int) string {

	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func randomInt(a, b int) int {

	var bytes = make([]byte, 1)
	rand.Read(bytes)

	per := float32(bytes[0]) / 256.0
	dif := maxInt(a, b) - minInt(a, b)

	return minInt(a, b) + int(per*float32(dif))
}

func SHA256(data []byte) []byte {
	sh := sha256.Sum256(data)
	return sh[:]
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
