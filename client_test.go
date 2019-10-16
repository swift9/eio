package eio_test

import (
	"github.com/swift9/eio"
	"testing"
	"time"
)

func test() {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	client := eio.NewClient("localhost:8000", protocol)
	var rpc *eio.RpcTemplate
	client.Connect(func(session *eio.Session) {
		rpc = eio.NewRpcTemplate(session)
		session.OnMessage = rpc.OnMessage
	})

	for i := 0; i < 1; i++ {
		go func() {
			for {
				m, _ := rpc.SendWithResponse(&eio.RpcMessage{
					MessageType:     []byte{0x00, 0x01},
					DataContentType: eio.TEXT,
					Body:            "hello",
				}, 1*time.Second)
				if m.ResponseId%10000 == 0 {
					println(time.Now().String(), m.RequestId)
				}
				if m.ResponseId > 30*10000 {
					break
				}
			}
		}()
	}
}

func test2() {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	client := eio.NewClient("localhost:8000", protocol)
	var rpc *eio.RpcTemplate
	client.Connect(func(session *eio.Session) {
		rpc = eio.NewRpcTemplate(session)
		rpc.RegisterRpcMessageHandle("0001", func(m *eio.RpcMessage) {
			if m.ResponseId%10000 == 0 {
				println(time.Now().String(), m.RequestId)
			}
		})
		session.OnMessage = rpc.OnMessage
	})

	for i := 0; i < 15; i++ {
		go func() {
			for {
				rpc.Send(&eio.RpcMessage{
					MessageType:     []byte{0x00, 0x01},
					DataContentType: eio.TEXT,
					Body:            "hello",
				}, 1*time.Second)
			}
		}()
	}
}

func TestClient_Rpc(t *testing.T) {
	go test2()
	time.Sleep(1 * time.Hour)
}
