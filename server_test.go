package eio_test

import (
	"github.com/swift9/eio"
	"log"
	"testing"
	"time"
)

type EchoProtocol struct {
}

func (p *EchoProtocol) Segment(buf *eio.ByteBuffer) []byte {
	if buf.Len() == 0 {
		return nil
	}
	bytes := buf.Read(0, int(buf.Len()))
	return bytes
}

func (p *EchoProtocol) Decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (p *EchoProtocol) Encode(d interface{}) ([]byte, error) {
	bytes, _ := d.([]byte)
	return bytes, nil
}

func (p *EchoProtocol) IsValidMessage(bytes []byte) bool {
	return true
}

func TestEcho_Server(t *testing.T) {
	server := eio.NewServer(":8000", &EchoProtocol{}, nil)

	go server.Listen(func(socket *eio.Session) {
		println(socket.Id)
		socket.OnSync("message", func(message interface{}) {
			bytes, _ := message.([]byte)
			socket.Write(bytes)
		})
	})

	client := eio.NewClient(":8000", &EchoProtocol{}, nil)

	client.Connect(func(s *eio.Session) {
		s.Write([]byte("hello"))
		s.On("message", func(message interface{}) {
			bytes, _ := message.([]byte)
			println("reply:" + string(bytes))
		})
	})
	time.Sleep(3 * time.Second)
}

type CustomProtocol struct {
	eio.VariableProtocol
}

func (p *CustomProtocol) Decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (p *CustomProtocol) Encode(d interface{}) ([]byte, error) {
	bytes, _ := d.([]byte)
	return bytes, nil
}

func (p *CustomProtocol) IsValidMessage(bytes []byte) bool {
	return true
}

func TestVariableProtocol_Server(t *testing.T) {
	customProtocol := &CustomProtocol{}
	customProtocol.MagicBytes = []byte{0xD0}
	customProtocol.MessageByteSize = 1
	server := eio.NewServer(":8000", customProtocol, nil)

	go server.Listen(func(socket *eio.Session) {
		println(socket.Id)
		socket.OnSync("message", func(message interface{}) {
			bytes, _ := message.([]byte)
			_, err := socket.Write(bytes)
			if err != nil {
				log.Println("write error", err)
			}
		})
	})

	time.Sleep(10 * time.Second)

}
