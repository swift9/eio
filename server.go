package eio

import (
	event "github.com/swift9/ares-event"
	"net"
)

type Server struct {
	event.Emitter
	tcpListener *net.TCPListener
	Sessions    map[string]*Session
	Addr        string
	Decoder     IDecodet
	Encoder     IEncoder
}

func NewServer(addr string) *Server {
	server := &Server{
		Addr: addr,
	}
	server.On("error", func(err error) {
		println(err)
	})
	return server
}

func (server *Server) Listen(onSession func(session *Session)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", server.Addr)

	if err != nil {
		return err
	}

	server.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := server.tcpListener.AcceptTCP()
			if err != nil {
				server.Emit("error", err)
				continue
			}
			session := NewSession(conn)
			server.Sessions[session.Id] = session
			onSession(session)
			go onSession(session)
		}
	}()
	return nil
}

func (server *Server) onSession(session *Session) {

}
