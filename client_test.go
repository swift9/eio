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
		if m.RequestId%10000 == 0 {
			println(time.Now().String(), s, m.RequestId)
		}
	})

	client.Connect(func(s *eio.Session) {
		go func() {
			i := 0
			println(time.Now().String())
			for {
				i++
				s.SendMessage(&eio.RpcMessage{
					MessageType:     []byte{0x00, 0x01},
					DataContentType: eio.TEXT,
					Body:            "hello",
				})
				if i > 266666 {
					break
				}
			}
		}()
	})

	time.Sleep(1 * time.Hour)
}

func TestClient_Rpc(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}

	rpc := eio.NewRpcTemplate()
	client := eio.NewClient("localhost:8000", protocol, rpc.OnMessage)

	client.Connect(rpc.OnConnect)

	m, err := rpc.SendWithResponse(&eio.RpcMessage{
		MessageType:     []byte{0x00, 0x01},
		DataContentType: eio.TEXT,
		Body:            "hello",
	}, 1*time.Second)

	if err == nil {
		println(m.ResponseId)
	}

	time.Sleep(1 * time.Hour)
}
