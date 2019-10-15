package eio

import (
	event "github.com/swift9/ares-event"
	"net"
)

type Server struct {
	event.Emitter
	tcpListener *net.TCPListener
	Sockets     map[string]*Session
	Addr        string
	Protocol    Protocol
	Log         ILog
	OnMessage   func(message interface{}, session *Session)
}

func NewServer(addr string, protocol Protocol, onMessage func(message interface{}, session *Session)) *Server {
	server := &Server{
		Addr:      addr,
		Protocol:  protocol,
		Sockets:   make(map[string]*Session),
		Log:       &SysLog{},
		OnMessage: onMessage,
	}
	return server
}

func (server *Server) SetOnMessage(f func(message interface{}, session *Session)) {
	server.OnMessage = f
}

func (server *Server) SetLog(log ILog) {
	server.Log = log
}

func (server *Server) Listen(onConnect func(session *Session)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", server.Addr)
	if err != nil {
		server.Log.Error(err)
		return err
	}

	server.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		server.Log.Error(err)
		return err
	}

	for {
		conn, err := server.tcpListener.AcceptTCP()
		if err != nil {
			server.Log.Error(err)
			server.Emit("error", err)
			continue
		}
		session := NewSession(conn, server.Protocol)
		session.SetLog(server.Log)

		if server.OnMessage != nil {
			session.On("message", func(message interface{}) {
				server.OnMessage(message, session)
			})
		}
		server.Sockets[session.Id] = session
		go func() {
			onConnect(session)
			session.poll()
		}()
	}
	return nil
}
