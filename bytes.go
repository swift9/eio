package eio

import "encoding/binary"

func Int8ToBytes(i int64) []byte {
	return []byte{byte(i)}
}

func Int16ToBytes(i int64) []byte {
	var buf = make([]byte, 16)
	binary.BigEndian.PutUint16(buf, uint16(i))
	return buf
}

func Int32ToBytes(i int64) []byte {
	var buf = make([]byte, 32)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 64)
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
