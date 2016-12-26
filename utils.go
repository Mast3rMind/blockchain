package blockchain

import "math/big"

func ZeroHash() []byte {
	return make([]byte, 32)
}

func concat(pieces ...[]byte) []byte {
	sz := 0
	for _, p := range pieces {
		sz += len(p)
	}

	buf := make([]byte, sz)

	i := 0
	for _, p := range pieces {
		copy(buf[i:], p)
		i += len(p)
	}
	return buf
}

func joinBigInt(expectedLen int, bigs ...*big.Int) *big.Int {
	bs := []byte{}
	for i, b := range bigs {
		by := b.Bytes()
		dif := expectedLen - len(by)
		if dif > 0 && i != 0 {
			by = append(arrayOfBytes(dif, 0), by...)
		}
		bs = append(bs, by...)
	}

	b := new(big.Int).SetBytes(bs)
	return b
}

func splitBigInt(b *big.Int, parts int) []*big.Int {
	bs := b.Bytes()
	if len(bs)%2 != 0 {
		bs = append([]byte{0}, bs...)
	}

	l := len(bs) / parts
	as := make([]*big.Int, parts)

	for i := range as {
		as[i] = new(big.Int).SetBytes(bs[i*l : (i+1)*l])
	}

	return as

}

// create an array filled with b
func arrayOfBytes(i int, b byte) (p []byte) {
	for i != 0 {
		p = append(p, b)
		i--
	}
	return
}

/*func powTargetFromBits(bits uint32) *big.Int {
	bits3 := bits - (bits>>24)<<24
	bits1 := 8 * (bits>>24 - 3)
	j := big.NewInt(int64(2))
	j.Exp(j, big.NewInt(int64(bits1)), big.NewInt(0))
	j.Mul(big.NewInt(int64(bits3)), j)
	return j
}*/
