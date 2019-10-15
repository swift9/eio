package eio

import (
	"bytes"
)

// 协议处理
// 报文分包、编解码
type Protocol interface {
	// 报文分包
	Segment(session *Session, buf *MessageByteBuffer) (start int64, end int64)

	// 解码
	Decode(session *Session, message *MessageByteBuffer) (interface{}, error)

	// 编码
	Encode(session *Session, d interface{}) (*MessageByteBuffer, error)

	// 验证报文有效性
	IsValidMessage(session *Session, message *MessageByteBuffer) bool
}

type VariableProtocol struct {
	MagicBytes      []byte
	MessageByteSize int
}

func (p *VariableProtocol) Segment(session *Session, messageByteBuffer *MessageByteBuffer) (start int64, end int64) {
	magicBytesLength := len(p.MagicBytes)

	headerLength := int64(magicBytesLength + p.MessageByteSize)
	if messageByteBuffer.Len() < headerLength {
		return 0, 0
	}

	if bytes.Equal(p.MagicBytes, messageByteBuffer.Peek(0, int64(magicBytesLength)).Message()) {
		messageByteBuffer.Discard(1)
		return p.Segment(session, messageByteBuffer)
	}

	lengthBytes := messageByteBuffer.buf[len(p.MagicBytes):(len(p.MagicBytes) + p.MessageByteSize)]
	messageLength := BytesToInt64(lengthBytes)

	if messageByteBuffer.Len() < messageLength {
		return 0, 0
	}

	return 0, messageLength
}

func (p *VariableProtocol) IsValidMessage(session *Session, message *MessageByteBuffer) bool {
	return true
}

func (p *VariableProtocol) Decode(session *Session, message *MessageByteBuffer) (interface{}, error) {
	return &message, nil
}

func (p *VariableProtocol) Encode(session *Session, d interface{}) (*MessageByteBuffer, error) {
	messageByte, _ := d.(MessageByteBuffer)
	return &messageByte, nil
}
