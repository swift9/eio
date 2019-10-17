package eio

import (
	"net"
)

type Server struct {
	tcpListener *net.TCPListener
	Addr        string
	Protocol    Protocol
	Log         ILog
}

func NewServer(addr string, protocol Protocol) *Server {
	server := &Server{
		Addr:     addr,
		Protocol: protocol,
		Log:      &SysLog{},
	}
	return server
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
			return err
		}
		session := NewSession(conn, server.Protocol)
		onConnect(session)
		if session.AutoPoll {
			session.poll()
		}
	}
	return nil
}
