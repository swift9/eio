package eio

import (
	"encoding/binary"
	"sync/atomic"
	"unsafe"
)

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	var i int64 = 0
	switch len(buf) {
	case 1:
		i = int64(buf[0])
	case 2:
		i = int64(binary.BigEndian.Uint16(buf))
	case 4:
		i = int64(binary.BigEndian.Uint32(buf))
	case 8:
		i = int64(binary.BigEndian.Uint64(buf))
	}
	return i
}

func Int642Int(i64 int64) int {
	i := (*int)(unsafe.Pointer(&i64))
	return *i
}

var seq int64 = 0

func GenerateSeq() int64 {
	return atomic.AddInt64(&seq, 1)
}
