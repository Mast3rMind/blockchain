package blockchain

import (
	"reflect"

	"github.com/ipkg/blockchain/utils"
)

const (
	KEY_POW_COMPLEXITY      = 0
	TEST_KEY_POW_COMPLEXITY = 0

	TRANSACTION_POW_COMPLEXITY      = 1
	TEST_TRANSACTION_POW_COMPLEXITY = 1

	BLOCK_POW_COMPLEXITY      = 2
	TEST_BLOCK_POW_COMPLEXITY = 2

	POW_PREFIX      = 0
	TEST_POW_PREFIX = 0
)

var (
	TRANSACTION_POW = utils.ArrayOfBytes(TRANSACTION_POW_COMPLEXITY, POW_PREFIX)
	BLOCK_POW       = utils.ArrayOfBytes(BLOCK_POW_COMPLEXITY, POW_PREFIX)

	TEST_TRANSACTION_POW = utils.ArrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, POW_PREFIX)
	TEST_BLOCK_POW       = utils.ArrayOfBytes(TEST_BLOCK_POW_COMPLEXITY, POW_PREFIX)
)

func CheckProofOfWork(prefix []byte, hash []byte) bool {
	if len(prefix) > 0 {
		return reflect.DeepEqual(prefix, hash[:len(prefix)])
	}
	return true
}
