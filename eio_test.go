package eio

import (
	"encoding/hex"
	"sync"
	"testing"
)

func Test_buffer(t *testing.T) {
	b1 := []byte{1, 2, 3}
	b2 := make([]byte, 3)
	copy(b2, b1)
	b2[0] = 3
	println(b1[0], b2[0])

}

func TestInt64ToBytes(t *testing.T) {
	println(hex.EncodeToString(Int64ToBytes(1)))
}

var lock = &sync.Mutex{}

func TestLock(t *testing.T) {
	test1()
}

func test1() {
	lock.Lock()
	defer lock.Unlock()
	println(1)
	test2()
}

func test2() {
	lock.Lock()
	defer lock.Unlock()
	println(2)
}
