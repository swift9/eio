package eio_test

import (
	"github.com/swift9/eio"
	"strconv"
	"testing"
	"time"
)

func TestServer_Listen(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	server := eio.NewServer(":8000", protocol, func(message interface{}, session *eio.Session) {
		mm, ok := (message).(*eio.RpcMessage)
		if ok {
			s, _ := mm.Body.(string)
			mm.Body = "reply:" + s + strconv.FormatInt(mm.RequestId, 10)
			session.SendMessage(mm)
		}
	})

	server.Listen(func(session *eio.Session) {
		println("connect", session.Conn.RemoteAddr().String())
	})
}

func TestServer_Rpc(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	rpc := eio.NewRpcTemplate()

	rpc.RegisterRpcMessageHandle("0001", func(message *eio.RpcMessage) {
		message.ResponseId = message.RequestId
		rpc.Send(message, 1*time.Second)
	})

	server := eio.NewServer(":8000", protocol, rpc.OnMessage)

	server.Listen(rpc.OnConnect)

	time.Sleep(1 * time.Hour)
}
