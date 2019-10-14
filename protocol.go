package eio

import (
	"encoding/hex"
)

// 协议处理
// 报文分包、编解码
type Protocol interface {
	// 报文分包
	Segment(buf *ByteBuffer) []byte

	// 解码
	Decode(bytes []byte) (interface{}, error)

	// 编码
	Encode(d interface{}) ([]byte, error)

	// 验证报文有效性
	IsValidMessage(message []byte) bool
}

type VariableProtocol struct {
	MagicBytes     []byte
	LengthByteSize int
}

func (p *VariableProtocol) Segment(buf *ByteBuffer) []byte {
	magicBytesLength := len(p.MagicBytes)

	headerLength := int64(magicBytesLength + p.LengthByteSize)
	if buf.Len() < headerLength {
		return nil
	}

	if hex.EncodeToString(p.MagicBytes) != hex.EncodeToString(buf.Read(0, magicBytesLength)) {
		buf.Discard(1)
		return p.Segment(buf)
	}

	lengthBytes := buf.buf[len(p.MagicBytes):(len(p.MagicBytes) + p.LengthByteSize)]
	messageLength := BytesToInt64(lengthBytes)

	if buf.Len() < messageLength {
		return nil
	}

	bytes := buf.Read(0, int(messageLength))
	return bytes
}

func (p *VariableProtocol) IsValidMessage(bytes []byte) bool {
	return true
}

func (p *VariableProtocol) Decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (p *VariableProtocol) Encode(d interface{}) ([]byte, error) {
	bytes, _ := d.([]byte)
	return bytes, nil
}
