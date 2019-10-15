package eio_test

import (
	"github.com/swift9/eio"
	"testing"
)

func TestClient_Connect(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}

	client := eio.NewClient("localhost:8000", protocol, func(message interface{}, session *eio.Session) {
		m, _ := message.(*eio.RpcMessage)
		println(m, m.Body)
	})

	client.Connect(func(s *eio.Session) {
		s.SendMessage(&eio.RpcMessage{
			MessageType: []byte{0x00, 0x01},
			Body:        "hello",
		})
	})
}
