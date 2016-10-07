package blockchain

import (
	"reflect"
)

// create an array filled with b
func arrayOfBytes(i int, b byte) (p []byte) {
	for i != 0 {
		p = append(p, b)
		i--
	}
	return
}

func fitBytesInto(d []byte, i int) []byte {
	if len(d) < i {
		dif := i - len(d)
		return append(arrayOfBytes(dif, 0), d...)
	}
	return d[:i]
}

// function map
// f = function, vs = slice
func funcMap(f interface{}, vs interface{}) interface{} {

	vf := reflect.ValueOf(f)
	vx := reflect.ValueOf(vs)

	l := vx.Len()

	tys := reflect.SliceOf(vf.Type().Out(0))
	vys := reflect.MakeSlice(tys, l, l)

	for i := 0; i < l; i++ {

		y := vf.Call([]reflect.Value{vx.Index(i)})[0]
		vys.Index(i).Set(y)
	}

	return vys.Interface()
}

func stripByte(d []byte, b byte) []byte {
	for i, bb := range d {
		if bb != b {
			return d[i:]
		}
	}
	return nil
}
