package eio

import (
	"bytes"
)

type VariableProtocol struct {
	MagicBytes      []byte
	MessageByteSize int
}

func (p *VariableProtocol) Segment(session *Session, messageByteBuffer *MessageBuffer) *MessageBuffer {
	magicBytesLength := len(p.MagicBytes)

	headerLength := magicBytesLength + p.MessageByteSize
	if messageByteBuffer.Len() < headerLength {
		return nil
	}

	if !bytes.Equal(p.MagicBytes, messageByteBuffer.Peek(0, magicBytesLength).Bytes()) {
		messageByteBuffer.Discard(1)
		return p.Segment(session, messageByteBuffer)
	}

	messageLength := messageByteBuffer.Int64Value(len(p.MagicBytes), len(p.MagicBytes)+p.MessageByteSize)

	if int64(messageByteBuffer.Len()) < messageLength {
		return nil
	}

	return messageByteBuffer.Peek(0, Int642Int(messageLength))
}

func (p *VariableProtocol) IsValidMessage(session *Session, message *MessageBuffer) bool {
	return true
}

func (p *VariableProtocol) Decode(session *Session, message *MessageBuffer) (interface{}, error) {
	return &message, nil
}

func (p *VariableProtocol) Encode(session *Session, d interface{}) (*MessageBuffer, error) {
	messageByte, _ := d.(MessageBuffer)
	return &messageByte, nil
}
