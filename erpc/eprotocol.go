package erpc

import (
	"errors"
	"github.com/swift9/eio"
)

// magicBytes MessageSize MessageType MessageId ResponseId DataType Data CheckCode
type EProtocol struct {
	eio.VariableProtocol
	MessageTypeByteSize int
	MessageIdByteSize   int
	CheckCodeByteSize   int
}

const (
	Text byte = 0xFF

	GzipText byte = 0x01

	Bin byte = 0x02

	GzipBin byte = 0x03
)

func NewDefaultEProtocol() *EProtocol {
	p := &EProtocol{
		VariableProtocol: eio.VariableProtocol{
			MagicBytes:            []byte{0xEE, 0xEE},
			MessageLengthByteSize: 8,
		},
		MessageTypeByteSize: 4,
		MessageIdByteSize:   8,
		CheckCodeByteSize:   0,
	}
	return p
}

func (rpcProtocol *EProtocol) IsValidMessage(session *eio.Session, message *eio.MessageBuffer) bool {
	if rpcProtocol.CheckCodeByteSize == 0 {
		return true
	}
	return true
}

func (rpcProtocol *EProtocol) GenerateCheckCode(b []byte) []byte {
	return nil
}

func (rpcProtocol *EProtocol) Decode(session *eio.Session, message *eio.MessageBuffer) (interface{}, error) {
	eioMessage := &EMessage{}
	MagicByteSize := len(rpcProtocol.MagicBytes)
	eioMessage.MessageSize = message.Int64Value(MagicByteSize, MagicByteSize+rpcProtocol.MessageIdByteSize)
	eioMessage.MessageType = message.Peek(MagicByteSize+rpcProtocol.MessageIdByteSize,
		MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize).Bytes()
	eioMessage.Id = message.Int64Value(MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize,
		MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+rpcProtocol.MessageIdByteSize)
	eioMessage.ResponseId = message.Int64Value(MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+rpcProtocol.MessageIdByteSize,
		MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+2*rpcProtocol.MessageIdByteSize)
	eioMessage.DataType = message.Peek(MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+2*rpcProtocol.MessageIdByteSize,
		MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+2*rpcProtocol.MessageIdByteSize+1).Bytes()[0]

	eioMessage.Data = message.Peek(MagicByteSize+rpcProtocol.MessageIdByteSize+rpcProtocol.MessageTypeByteSize+2*rpcProtocol.MessageIdByteSize+1,
		eio.Int642Int(eioMessage.MessageSize)-rpcProtocol.CheckCodeByteSize).Bytes()

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
			break
		}
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
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
	MagicByteSize := len(rpcProtocol.MagicBytes)
	size := MagicByteSize + rpcProtocol.MessageIdByteSize + rpcProtocol.MessageTypeByteSize + 2*rpcProtocol.MessageIdByteSize + 1 + dataLength + rpcProtocol.CheckCodeByteSize
	byteBuffer.Append(eio.Int64ToBytes(int64(size)))
	byteBuffer.Append(rpcMessage.MessageType)
	byteBuffer.Append(eio.Int64ToBytes(rpcMessage.Id))
	byteBuffer.Append(eio.Int64ToBytes(rpcMessage.ResponseId))
	byteBuffer.AppendByte(rpcMessage.DataType)
	byteBuffer.Append(data)

	if code := rpcProtocol.GenerateCheckCode(byteBuffer.Bytes()); code != nil {
		byteBuffer.Append(code)
	}

	return byteBuffer, nil
}
