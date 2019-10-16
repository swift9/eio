package eio_test

import (
	"github.com/swift9/eio"
	"testing"
	"time"
)

func TestClient_Connect(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}

	client := eio.NewClient("localhost:8000", protocol, func(message interface{}, session *eio.Session) {
		m, _ := message.(*eio.RpcMessage)
		s, _ := m.Body.(string)
		println(s)
	})

	client.Connect(func(s *eio.Session) {
		s.SendMessage(&eio.RpcMessage{
			MessageType:     []byte{0x00, 0x01},
			DataContentType: eio.TEXT,
			Body:            "hello",
		})
	})

	time.Sleep(1 * time.Hour)
}
