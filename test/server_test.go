package eio_test

import (
	"github.com/swift9/eio/erpc"
	"testing"
	"time"
)

func TestServer_Rpc(t *testing.T) {
	server := erpc.NewEServer(":8000", erpc.NewDefaultEProtocol())
	server.Listen(func(session *erpc.ESession) {
		session.RegisterMessageHandle("00000001", func(message *erpc.EMessage) {
			message.ResponseId = message.Id
			session.Send(message, 1*time.Second)
		})
	})
}
