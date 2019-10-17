package erpc

import (
	"errors"
	"github.com/swift9/eio"
)

type EProtocol struct {
	eio.VariableProtocol
	CheckCodeBytes []byte
}

const (
	Text byte = 0xFF

	GzipText byte = 0x01

	Bin byte = 0x02

	GzipBin byte = 0x03
)

func (rpcProtocol *EProtocol) GetCheckCodeLength() int {
	if rpcProtocol.CheckCodeBytes == nil {
		return 0
	}
	return len(rpcProtocol.CheckCodeBytes)
}

func (rpcProtocol *EProtocol) IsValidMessage(session *eio.Session, message *eio.MessageBuffer) bool {
	if rpcProtocol.CheckCodeBytes == nil {
		return true
	}
	return true
}

func (rpcProtocol *EProtocol) Decode(session *eio.Session, message *eio.MessageBuffer) (interface{}, error) {
	eioMessage := &EMessage{}
	eioMessage.MessageSize = message.Int64Value(2, 10)
	eioMessage.MessageType = message.Peek(10, 12).Bytes()
	eioMessage.Id = message.Int64Value(12, 20)
	eioMessage.ResponseId = message.Int64Value(20, 28)
	eioMessage.DataType = message.Peek(28, 29).Bytes()[0]

	eioMessage.Data = message.Peek(29, eio.Int642Int(eioMessage.MessageSize)-rpcProtocol.GetCheckCodeLength()).Bytes()
	return eioMessage, nil
}

func (rpcProtocol *EProtocol) Encode(session *eio.Session, message interface{}) (*eio.MessageBuffer, error) {
	byteBuffer := eio.NewMessageBuffer()
	byteBuffer.Append(rpcProtocol.MagicBytes)

	rpcMessage, _ := message.(*EMessage)
	var data []byte
	dataLength := 0

	switch rpcMessage.DataType {
	case Text:
		if d, ok := rpcMessage.Data.(string); ok {
			data = []byte(d)
			dataLength = len(d)
		}
	case GzipText:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	case Bin:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	case GzipBin:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	default:
		return nil, errors.New("not support")
	}
	byteBuffer.Append(eio.Int64ToBytes(29 + int64(dataLength) + int64(rpcProtocol.GetCheckCodeLength())))
	byteBuffer.Append(rpcMessage.MessageType)
	byteBuffer.Append(eio.Int64ToBytes(rpcMessage.Id))
	byteBuffer.Append(eio.Int64ToBytes(rpcMessage.ResponseId))
	byteBuffer.AppendByte(rpcMessage.DataType)
	byteBuffer.Append(data)
	if rpcProtocol.CheckCodeBytes != nil {
		byteBuffer.Append(rpcProtocol.CheckCodeBytes)
	}

	return byteBuffer, nil
}
