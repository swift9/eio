package eio

import (
	"bufio"
	"errors"
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
	Log             ILog
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
		Log:             &SysLog{},
	}
	return socket
}

func (s *Socket) SetLog(log ILog) {
	s.Log = log
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
	if s.Conn == nil {
		return 0, errors.New("connection is closed")
	}

	return s.Conn.Write(bytes)
}

func (s *Socket) WriteData(data interface{}) (int, error) {
	bytes, _ := s.Protocol.Encode(data)
	return s.Write(bytes)
}

func (s *Socket) Read(bytes []byte) (int, error) {
	if s.Conn == nil {
		s.Log.Error("connection is nil")
		return 0, errors.New("connection is closed")
	}
	return s.Conn.Read(bytes)
}

func (s *Socket) CloseRead() error {
	return s.Conn.CloseRead()
}

func (s *Socket) CloseWrite() error {
	return s.Conn.CloseWrite()
}

func (s *Socket) Close() error {
	e := s.Conn.Close()
	return e
}

func (s *Socket) Pipe(socket *Socket) {
	w := bufio.NewWriterSize(socket, socket.WriteBufferSize)
	r := bufio.NewReaderSize(s, s.ReadBufferSize)
	w.ReadFrom(r)
}

func (s *Socket) byteBufferPoll() (int, error) {
	var (
		n         = 0
		err error = nil
	)
	bytes := make([]byte, s.ReadBufferSize)
	if n, err = s.Conn.Read(bytes); err != nil {
		s.Log.Error("socket read error ", err)
		s.Emit("error", err)
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
	defer func() {
		if err := recover(); err != nil {
			socket.Log.Error("poll ", err)
			socket.Emit("error", errors.New("unknown"))
		}
	}()
	for {
		n, err := socket.byteBufferPoll()
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
			if !s.Protocol.IsValidMessage(bytes) {
				s.ByteBuffer.Discard(1)
				s.ByteBufferSegment()
				return
			} else {
				s.ByteBuffer.Discard(len(bytes))
			}
			if data, err := s.Protocol.Decode(bytes); err == nil {
				s.Emit("data", data)
			} else {
				s.Log.Error("decode ", err)
			}
		} else {
			return
		}
	}
}
