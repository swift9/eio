package eio

import (
	"bufio"
	"errors"
	"net"
)

type Session struct {
	Id                 int64
	ReadBufferSize     int
	WriteBufferSize    int
	MessageByteBuffer  *MessageBuffer
	Conn               *net.TCPConn
	IsClosed           bool
	Protocol           Protocol
	isPooled           bool
	Log                ILog
	Context            map[string]interface{}
	OnMessage          func(message interface{}, session *Session)
	OnReadOrWriteError func(err error, session *Session)
}

func NewSession(conn *net.TCPConn, protocol Protocol) *Session {
	socket := &Session{
		Id:                GenerateSeq(),
		Conn:              conn,
		ReadBufferSize:    1024 * 4,
		WriteBufferSize:   1024 * 256,
		MessageByteBuffer: NewMessageByteBuffer(),
		Protocol:          protocol,
		isPooled:          false,
		Log:               &SysLog{},
		OnReadOrWriteError: func(err error, session *Session) {
			session.Close()
		},
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

func (s *Session) SendMessage(message interface{}) (int, error) {
	bytes, _ := s.Protocol.Encode(s, message)
	return s.Write(bytes.Bytes())
}

func (s *Session) Write(bytes []byte) (int, error) {

	if s.IsClosed || s.Conn == nil {
		return 0, errors.New("connection is closed")
	}
	n, err := s.Conn.Write(bytes)
	if err != nil {
		s.Log.Error("write error ", err)
		s.OnReadOrWriteError(err, s)
	}
	return n, err
}

func (s *Session) Read(bytes []byte) (int, error) {
	if s.IsClosed || s.Conn == nil {
		s.Log.Error("connection is closed")
		return 0, errors.New("connection is closed")
	}
	n, err := s.Conn.Read(bytes)
	if err != nil {
		s.Log.Error("read error ", err)
		s.OnReadOrWriteError(err, s)
	}
	return n, err
}

func (s *Session) Close() error {
	if s.IsClosed {
		return nil
	}
	s.IsClosed = true
	e := s.Conn.Close()
	return e
}

func (s *Session) Pipe(session *Session) {
	w := bufio.NewWriterSize(session, session.WriteBufferSize)
	r := bufio.NewReaderSize(s, s.ReadBufferSize)
	w.ReadFrom(r)
}

func (s *Session) poll() {
	if s.isPooled {
		return
	} else {
		s.isPooled = true
	}
	go func() {
		var (
			n         = 0
			err error = nil
		)
		for {
			bytes := make([]byte, s.ReadBufferSize)
			if n, err = s.Read(bytes); err != nil {
				return
			}
			if n > 0 {
				s.MessageByteBuffer.Append(bytes[0:n])
				s.segment()
			}
		}
	}()
}

func (s *Session) segment() {
	for {
		messageByteBuffer := s.Protocol.Segment(s, s.MessageByteBuffer)

		if messageByteBuffer == nil {
			return
		}
		if !s.Protocol.IsValidMessage(s, messageByteBuffer) {
			s.MessageByteBuffer.Discard(1)
			continue
		}
		s.MessageByteBuffer.Discard(messageByteBuffer.Len())
		go s.handleMessageByteBuffer(messageByteBuffer)
	}
}

func (s *Session) handleMessageByteBuffer(messageByte *MessageBuffer) {
	message, err := s.Protocol.Decode(s, messageByte)
	if err != nil {
		s.Log.Error("decode ", err)
	}
	go s.OnMessage(message, s)
}
