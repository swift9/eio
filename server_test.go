package eio_test

import (
	"eio"
	"testing"
)

type EchoProtocol struct {
}

func (p *EchoProtocol) Segment(buf *eio.ByteBuffer) []byte {
	if buf.Len() == 0 {
		return nil
	}
	bytes := buf.Read(0, int(buf.Len()))
	buf.Discard(int(buf.Len()))
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

func TestServer(t *testing.T) {
	server := eio.NewServer(":8000", &EchoProtocol{})

	server.Listen(func(socket *eio.Socket) {
		println(socket.Id)
		socket.OnSync("data", func(data interface{}) {
			bytes, _ := data.([]byte)
			socket.Write(bytes)
		})
		socket.Poll()
	})

}
