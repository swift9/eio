package eio

import "testing"

func Test_buffer(t *testing.T) {
	b1 := []byte{1, 2, 3}
	b2 := make([]byte, 3)
	copy(b2, b1)
	b2[0] = 3
	println(b1[0], b2[0])

}
