package eio

import (
	uuid "github.com/satori/go.uuid"
	event "github.com/swift9/ares-event"
	"net"
	"reflect"
)

type Session struct {
	event.Emitter
	Id              string
	Conn            *net.TCPConn
	decoder         IDecodet
	encoder         IEncoder
	readBufferSize  int
	writeBufferSize int
	tcpNoDelay      bool
}

func NewSession(conn *net.TCPConn) *Session {
	id := uuid.NewV4().String()
	session := &Session{
		Id:              id,
		Conn:            conn,
		readBufferSize:  1024 * 16,
		writeBufferSize: 1024 * 16,
		tcpNoDelay:      false,
	}
	return session
}

func (session *Session) SetReadBufferSize(size int) {
	session.Conn.SetReadBuffer(size)
}
func (session *Session) SetWriteBufferSize(size int) {
	session.Conn.SetWriteBuffer(size)
}
func (session *Session) SetTcpNoDelay(tcpNoDelay bool) {
	session.Conn.SetNoDelay(tcpNoDelay)
}

func (session *Session) Write(data interface{}) error {
	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.Slice:
		if v, ok := data.([]byte); ok {
			_, err := session.Conn.Write(v)
			return err
		}
	case reflect.String:
		if v, ok := data.(string); ok {
			_, err := session.Conn.Write([]byte(v))
			return err
		}
	default:
		bytes, err := session.encoder.Encode(data)
		if err != nil {
			return err
		}
		_, err = session.Conn.Write(bytes)
		return err
	}
	return nil
}

func (session *Session) Read() (data interface{}, err error) {
	return nil, nil
}
