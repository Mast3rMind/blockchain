package utils

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -(MaxInt - 1)
)

func ArrayOfBytes(i int, b byte) (p []byte) {
	for i != 0 {
		p = append(p, b)
		i--
	}
	return
}

func FitBytesInto(d []byte, i int) []byte {
	if len(d) < i {
		dif := i - len(d)
		return append(ArrayOfBytes(dif, 0), d...)
	}
	return d[:i]
}

func StripByte(d []byte, b byte) []byte {
	for i, bb := range d {
		if bb != b {
			return d[i:]
		}
	}
	return nil
}
