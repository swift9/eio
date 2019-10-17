package eio_test

import (
	"github.com/swift9/eio"
	"testing"
	"time"
)

func TestServer_Rpc(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}

	server := eio.NewServer(":8000", protocol)

	server.Listen(func(session *eio.Session) {
		rpc := eio.NewRpcTemplate(session)
		session.OnMessage = rpc.OnMessage

		rpc.RegisterRpcMessageHandle("0001", func(message *eio.RpcMessage) {
			message.ResponseId = message.RequestId
			rpc.Send(message, 1*time.Second)
		})
	})

	for {
		time.Sleep(1 * time.Hour)
	}
}
