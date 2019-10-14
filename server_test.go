package eio_test

import (
	"eio"
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
	server := eio.NewServer(":8000", &EchoProtocol{})

	go server.Listen(func(socket *eio.Socket) {
		println(socket.Id)
		socket.OnSync("data", func(data interface{}) {
			bytes, _ := data.([]byte)
			socket.Write(bytes)
		})
		socket.Poll()
	})

	client := eio.NewClient(":8000", &EchoProtocol{})

	client.Connect(func(s *eio.Socket) {
		s.Poll()
		s.Write([]byte("hello"))
		s.On("data", func(data interface{}) {
			bytes, _ := data.([]byte)
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
	customProtocol.LengthByteSize = 1
	server := eio.NewServer(":8000", customProtocol)

	go server.Listen(func(socket *eio.Socket) {
		println(socket.Id)
		socket.OnSync("data", func(data interface{}) {
			bytes, _ := data.([]byte)
			_, err := socket.Write(bytes)
			if err != nil {
				log.Println("write error", err)
			}
		})
		socket.Poll()
	})

	time.Sleep(10 * time.Second)

}
