package eio

import "bytes"

// 协议处理
// 报文分包、编解码
type Protocol interface {
	// 报文分包
	Segment(session *Session, buf *MessageByteBuffer) (start int, end int)

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

func (p *VariableProtocol) Segment(session *Session, messageByteBuffer *MessageByteBuffer) (start int, end int) {
	magicBytesLength := len(p.MagicBytes)

	headerLength := magicBytesLength + p.MessageByteSize
	if messageByteBuffer.Len() < headerLength {
		return 0, 0
	}

	if !bytes.Equal(p.MagicBytes, messageByteBuffer.Peek(0, magicBytesLength).Message()) {
		messageByteBuffer.Discard(1)
		return p.Segment(session, messageByteBuffer)
	}

	messageLength := messageByteBuffer.Int64Value(len(p.MagicBytes), len(p.MagicBytes)+p.MessageByteSize)

	if int64(messageByteBuffer.Len()) < messageLength {
		return 0, 0
	}

	return 0, Int642Int(messageLength)
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
