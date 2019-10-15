package eio

import "errors"

type RpcProtocol struct {
	VariableProtocol
	CheckCodeBytes []byte
}

const (
	TEXT byte = 0xFF

	GZIP_TEXT byte = 0x01

	BIN byte = 0x02

	GZIP_BIN byte = 0x03
)

func (rpcProtocol *RpcProtocol) GetCheckCodeLength() int64 {
	if rpcProtocol.CheckCodeBytes == nil {
		return 0
	}
	return int64(len(rpcProtocol.CheckCodeBytes))
}

func (rpcProtocol *RpcProtocol) IsValidMessage(session *Session, message *MessageByteBuffer) bool {
	if rpcProtocol.CheckCodeBytes == nil {
		return true
	}
	return true
}

func (rpcProtocol *RpcProtocol) Decode(session *Session, message *MessageByteBuffer) (interface{}, error) {
	eioMessage := &RpcMessage{}
	eioMessage.MessageSize = message.Int64Value(2, 9)
	eioMessage.MessageType = message.Peek(10, 11).Message()
	eioMessage.RequestId = message.Int64Value(12, 19)
	eioMessage.ResponseId = message.Int64Value(20, 27)
	eioMessage.DataContentType = message.Peek(28, 28).Message()[0]
	if eioMessage.DataContentType == TEXT {
		eioMessage.Body = string(message.Peek(29, eioMessage.MessageSize-1-rpcProtocol.GetCheckCodeLength()).Message())
	}
	return eioMessage, nil
}

var requestIndex int64

func (rpcProtocol *RpcProtocol) Encode(session *Session, message interface{}) (*MessageByteBuffer, error) {
	byteBuffer := &MessageByteBuffer{}
	byteBuffer.Append(rpcProtocol.MagicBytes)

	rpcMessage, _ := message.(*RpcMessage)
	if str, ok := rpcMessage.Body.(string); ok {
		if rpcMessage.DataContentType == TEXT {
			body := []byte(str)
			byteBuffer.Append(Int64ToBytes(29 + int64(len(body)) + rpcProtocol.GetCheckCodeLength()))
			byteBuffer.Append(rpcMessage.MessageType)
			requestIndex++
			byteBuffer.Append(Int64ToBytes(requestIndex))
			byteBuffer.Append(Int64ToBytes(rpcMessage.ResponseId))
			byteBuffer.Append(Int64ToBytes(rpcMessage.ResponseId))
			byteBuffer.AppendByte(rpcMessage.DataContentType)
			byteBuffer.Append(body)
			if rpcProtocol.CheckCodeBytes != nil {
				byteBuffer.Append(rpcProtocol.CheckCodeBytes)
			}
		}
	} else {
		return nil, errors.New("not support")
	}
	return byteBuffer, nil
}

func NewRpcProtocol() {

}
