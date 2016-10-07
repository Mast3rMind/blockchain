package blockchain

import (
	"reflect"
	"testing"
)

func Test_Transaction_Marshalling(t *testing.T) {

	kp := GenerateNewKeypair()
	tr := NewTransaction(kp.Public, nil, []byte(randomString(randomInt(0, 1024*1024))))

	tr.Header.Nonce = tr.GenerateNonce(arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX))
	tr.Signature = tr.Sign(kp)

	data, err := tr.MarshalBinary()

	if err != nil {
		t.Error(err)
	}

	newT := &Transaction{}
	rem, err := newT.UnmarshalBinary(data)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*newT, *tr) || len(rem) < 0 {
		t.Error("Marshall, unmarshall failed")
	}
}

func Test_Transaction_VerifyTransaction(t *testing.T) {

	pow := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX)

	kp := GenerateNewKeypair()
	tr := NewTransaction(kp.Public, nil, []byte(randomString(randomInt(0, 1024))))

	tr.Header.Nonce = tr.GenerateNonce(pow)
	tr.Signature = tr.Sign(kp)

	if !tr.VerifyTransaction(pow) {

		t.Error("Validation failing")
	}
}

func TestIncorrectTransactionPOWVerification(t *testing.T) {

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

func TestIncorrectTransactionSignatureVerification(t *testing.T) {

	pow := arrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, TEST_POW_PREFIX)
	kp1, kp2 := GenerateNewKeypair(), GenerateNewKeypair()
	tr := NewTransaction(kp2.Public, nil, []byte(randomString(randomInt(0, 1024))))
	tr.Header.Nonce = tr.GenerateNonce(pow)
	tr.Signature = tr.Sign(kp1)

	if tr.VerifyTransaction(pow) {

		t.Error("Passed validation with incorrect key")
	}
}
