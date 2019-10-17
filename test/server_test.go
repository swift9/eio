package eio_test

import (
	"github.com/swift9/eio/erpc"
	"testing"
	"time"
)

func TestServer_Rpc(t *testing.T) {
	protocol := &erpc.EProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}
	server := erpc.NewEServer(":8000", protocol)
	server.Listen(func(session *erpc.ESession) {
		session.RegisterMessageHandle("0001", func(message *erpc.EMessage) {
			message.ResponseId = message.Id
			session.Send(message, 1*time.Second)
		})
	})
}
