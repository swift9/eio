package eio

import (
	event "github.com/swift9/ares-event"
	"net"
)

type Client struct {
	event.Emitter
	Addr      string
	Protocol  Protocol
	Log       ILog
	OnMessage func(message interface{}, session *Session)
}

func NewClient(addr string, protocol Protocol, onMessage func(message interface{}, session *Session)) *Client {
	client := &Client{
		Addr:      addr,
		Protocol:  protocol,
		Log:       &SysLog{},
		OnMessage: onMessage,
	}
	return client
}

func (c *Client) SetLog(log ILog) {
	c.Log = log
}

func (c *Client) Connect(onConnect func(s *Session)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.Log.Error("connection ", err)
		return err
	}
	session := NewSession(conn, c.Protocol)
	session.SetLog(c.Log)
	if c.OnMessage != nil {
		session.On("message", func(message interface{}) {
			c.OnMessage(message, session)
		})
	}
	go func() {
		onConnect(session)
		session.poll()
	}()
	return nil
}
