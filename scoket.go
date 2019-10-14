package eio

import (
	"bufio"
	uuid "github.com/satori/go.uuid"
	event "github.com/swift9/ares-event"
	"net"
)

type Socket struct {
	event.Emitter
	Id              string
	Conn            *net.TCPConn
	ReadBufferSize  int
	WriteBufferSize int
	TcpNoDelay      bool
	ByteBuffer      *ByteBuffer
	Protocol        Protocol
	isPooled        bool
}

func NewSocket(conn *net.TCPConn, protocol Protocol) *Socket {
	id := uuid.NewV4().String()
	socket := &Socket{
		Id:              id,
		Conn:            conn,
		ReadBufferSize:  1024 * 1024,
		WriteBufferSize: 1024 * 1024,
		TcpNoDelay:      false,
		ByteBuffer:      &ByteBuffer{},
		Protocol:        protocol,
		isPooled:        false,
	}
	return socket
}

func (s *Socket) SetReadBufferSize(size int) {
	s.Conn.SetReadBuffer(size)
}

func (s *Socket) SetWriteBufferSize(size int) {
	s.Conn.SetWriteBuffer(size)
}

func (s *Socket) SetTcpNoDelay(tcoNoDelay bool) {
	s.Conn.SetNoDelay(tcoNoDelay)
}

func (s *Socket) Write(bytes []byte) (int, error) {
	return s.Conn.Write(bytes)
}

func (s *Socket) Read(bytes []byte) (int, error) {
	return s.Conn.Read(bytes)
}

func (s *Socket) Close(err error) error {
	e := s.Conn.Close()
	s.Emit("close", err)
	return e
}

func (s *Socket) Pipe(socket *Socket) {
	w := bufio.NewWriterSize(socket, socket.WriteBufferSize)
	r := bufio.NewReaderSize(s, s.ReadBufferSize)
	w.ReadFrom(r)
}

func (s *Socket) ByteBufferPoll() (int, error) {
	var (
		n         = 0
		err error = nil
	)
	bytes := make([]byte, s.ReadBufferSize)
	if n, err = s.Conn.Read(bytes); err != nil {
		s.Close(err)
		return 0, err
	}
	s.ByteBuffer.Append(bytes[0:n])
	return n, nil
}

func (s *Socket) Poll() {
	if !s.isPooled {
		s.isPooled = true
		go poll(s)
	}
}

func poll(socket *Socket) {
	for {
		n, err := socket.ByteBufferPoll()
		if err != nil {
			break
		}
		if n > 0 {
			socket.ByteBufferSegment()
		}
	}
}

func (s *Socket) ByteBufferSegment() {
	for {
		if bytes := s.Protocol.Segment(s.ByteBuffer); bytes != nil {
			if data, err := s.Protocol.Decode(bytes); err == nil {
				b, _ := data.([]byte)
				s.Emit("data", b)
			}
		} else {
			break
		}
	}
}
