package eio

import (
	event "github.com/swift9/ares-event"
	"net"
)

type Client struct {
	event.Emitter
	Addr     string
	Protocol Protocol
	Log      ILog
}

func NewClient(addr string, protocol Protocol) *Client {
	client := &Client{
		Addr:     addr,
		Protocol: protocol,
		Log:      &SysLog{},
	}
	return client
}

func (c *Client) SetLog(log ILog) {
	c.Log = log
}

func (c *Client) Connect(onConnect func(s *Socket)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Addr)

	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.Log.Error("connection ", err)
		return err
	}
	socket := NewSocket(conn, c.Protocol)
	socket.SetLog(c.Log)
	go onConnect(socket)
	return nil
}
