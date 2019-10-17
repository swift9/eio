package erpc

import "github.com/swift9/eio"

type EServer struct {
	Server *eio.Server
}

func NewEServer(addr string, protocol *EProtocol) *EServer {
	server := eio.NewServer(addr, protocol)
	return &EServer{
		Server: server,
	}
}

func (server *EServer) Listen(onConnect func(es *ESession)) {
	server.Server.Listen(func(s *eio.Session) {
		es := NewESession(s)
		s.OnMessage = es.onMessage
		onConnect(es)
	})
}
