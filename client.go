package eio

import (
	event "github.com/swift9/ares-event"
	"log"
	"net"
)

type Client struct {
	event.Emitter
	Addr     string
	Protocol Protocol
}

func NewClient(addr string, protocol Protocol) *Client {
	client := &Client{
		Addr:     addr,
		Protocol: protocol,
	}
	return client
}

func (c *Client) Connect(onConnect func(s *Socket)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Addr)

	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("ERROR", err)
		return err
	}
	socket := NewSocket(conn, c.Protocol)
	go onConnect(socket)
	return nil
}
