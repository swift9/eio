package eio_test

import (
	"encoding/hex"
	"github.com/swift9/eio"
	"testing"
	"time"
)

type CustomProtocolClient struct {
	eio.VariableProtocol
}

func (p *CustomProtocolClient) Decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (p *CustomProtocolClient) Encode(d interface{}) ([]byte, error) {
	bytes, _ := d.([]byte)
	return bytes, nil
}

func (p *CustomProtocolClient) IsValidMessage(bytes []byte) bool {
	return true
}

func TestVariableProtocol_Client(t *testing.T) {
	customProtocol := &CustomProtocolClient{}
	customProtocol.MagicBytes = []byte{0xD0}
	customProtocol.MessageByteSize = 1

	client := eio.NewClient(":8000", customProtocol, nil)

	client.Connect(func(s *eio.Session) {
		s.Write([]byte{0xD0, 0x03, 0x00, 0xD0, 0x03, 0x88})
		s.On("message", func(message interface{}) {
			bytes, _ := message.([]byte)
			println("server reply:" + hex.EncodeToString(bytes))
		})
	})
	time.Sleep(10 * time.Second)
}
