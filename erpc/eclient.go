package erpc

import (
	"github.com/swift9/eio"
)

type EClient struct {
	Client *eio.Client
}

func NewEClient(addr string, protocol *EProtocol) *EClient {
	client := &eio.Client{
		Addr:     addr,
		Protocol: protocol,
		Log:      &eio.SysLog{},
	}
	return &EClient{Client: client}
}

func (ec *EClient) Connect(onConnect func(es *ESession)) error {
	err := ec.Client.Connect(func(s *eio.Session) {
		es := NewESession(s)
		s.OnMessage = es.onMessage
		s.OnReadOrWriteError = func(err error, session *eio.Session) {
			es.Close()
		}
		onConnect(es)
	})
	return err
}
