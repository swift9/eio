package eio_test

import (
	"github.com/swift9/eio"
	"testing"
)

func TestServer_Listen(t *testing.T) {
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	server := eio.NewServer(":8000", protocol, func(message interface{}, session *eio.Session) {
		m, _ := message.(*eio.RpcMessage)
		println(m, m.Body)
		session.SendMessage(m)
	})

	server.Listen(func(session *eio.Session) {
		println("connect", session.Conn.RemoteAddr().String())
	})
}
