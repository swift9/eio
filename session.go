package eio

import (
	"bufio"
	"errors"
	uuid "github.com/satori/go.uuid"
	event "github.com/swift9/ares-event"
	"net"
)

type Session struct {
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
	isClosedRead    bool
	isClosedWrite   bool
}

func NewSession(conn *net.TCPConn, protocol Protocol) *Session {
	id := uuid.NewV4().String()
	socket := &Session{
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

func (s *Session) SetLog(log ILog) {
	s.Log = log
}

func (s *Session) SetReadBufferSize(size int) {
	s.Conn.SetReadBuffer(size)
}

func (s *Session) SetWriteBufferSize(size int) {
	s.Conn.SetWriteBuffer(size)
}

func (s *Session) SetTcpNoDelay(tcoNoDelay bool) {
	s.Conn.SetNoDelay(tcoNoDelay)
}

func (s *Session) SetKeepAlive(keepAlive bool) {
	s.Conn.SetKeepAlive(keepAlive)
}

func (s *Session) Write(bytes []byte) (int, error) {
	if s.Conn == nil {
		return 0, errors.New("connection is closed")
	}
	n, err := s.Conn.Write(bytes)

	if err != nil {
		s.Log.Error("write error ", err)
		s.CloseWrite()
	}
	return n, err
}

func (s *Session) SendMessage(message interface{}) (int, error) {
	bytes, _ := s.Protocol.Encode(message)
	return s.Write(bytes)
}

func (s *Session) Read(bytes []byte) (int, error) {
	if s.Conn == nil {
		s.Log.Error("connection is nil")
		return 0, errors.New("connection is closed")
	}
	n, err := s.Conn.Read(bytes)

	if err != nil {
		s.Log.Error("read error ", err)
		s.CloseRead()
	}
	return n, err
}

func (s *Session) CloseRead() error {
	err := s.Conn.CloseRead()
	s.isClosedRead = true
	s.Emit("closeRead")
	return err
}

func (s *Session) CloseWrite() error {
	err := s.Conn.CloseWrite()
	s.isClosedWrite = true
	s.Emit("closeWrite")
	return err
}

func (s *Session) Close() error {
	e := s.Conn.Close()
	s.Emit("close")
	return e
}

func (s *Session) Pipe(socket *Session) {
	w := bufio.NewWriterSize(socket, socket.WriteBufferSize)
	r := bufio.NewReaderSize(s, s.ReadBufferSize)
	w.ReadFrom(r)
}

func (s *Session) readIn() (int, error) {
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

func (s *Session) poll() {
	if !s.isPooled {
		s.isPooled = true
		go poll(s)
	}
}

func poll(session *Session) {
	defer func() {
		if err := recover(); err != nil {
			session.Log.Error("poll ", err)
			session.Emit("error", errors.New("unknown"))
			session.CloseRead()
		}
	}()
	for {
		n, err := session.readIn()
		if err != nil {
			break
		}
		if n > 0 {
			session.ByteBufferSegment()
		}
	}
}

func (s *Session) ByteBufferSegment() {
	for {
		if bytes := s.Protocol.Segment(s.ByteBuffer); bytes != nil {
			if !s.Protocol.IsValidMessage(bytes) {
				s.ByteBuffer.Discard(1)
				s.ByteBufferSegment()
				return
			} else {
				s.ByteBuffer.Discard(len(bytes))
			}
			if message, err := s.Protocol.Decode(bytes); err == nil {
				s.Emit("message", message)
			} else {
				s.Log.Error("decode ", err)
			}
		} else {
			return
		}
	}
}
