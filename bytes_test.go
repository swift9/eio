package eio

import (
	"encoding/hex"
	"testing"
)

func TestInt64ToBytes(t *testing.T) {
	println(hex.EncodeToString(Int64ToBytes(1)))
}
