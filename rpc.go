package eio

import (
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type RpcMessage struct {
	Id          int64
	ResponseId  int64
	MessageSize int64
	MessageType []byte
	DataType    byte
	Data        interface{}
}

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

func (rpcProtocol *RpcProtocol) GetCheckCodeLength() int {
	if rpcProtocol.CheckCodeBytes == nil {
		return 0
	}
	return len(rpcProtocol.CheckCodeBytes)
}

func (rpcProtocol *RpcProtocol) IsValidMessage(session *Session, message *MessageByteBuffer) bool {
	if rpcProtocol.CheckCodeBytes == nil {
		return true
	}
	return true
}

func (rpcProtocol *RpcProtocol) Decode(session *Session, message *MessageByteBuffer) (interface{}, error) {
	eioMessage := &RpcMessage{}
	eioMessage.MessageSize = message.Int64Value(2, 10)
	eioMessage.MessageType = message.Peek(10, 12).Message()
	eioMessage.Id = message.Int64Value(12, 20)
	eioMessage.ResponseId = message.Int64Value(20, 28)
	eioMessage.DataType = message.Peek(28, 29).Message()[0]

	eioMessage.Data = message.Peek(29, Int642Int(eioMessage.MessageSize)-rpcProtocol.GetCheckCodeLength()).Message()
	return eioMessage, nil
}

func (rpcProtocol *RpcProtocol) Encode(session *Session, message interface{}) (*MessageByteBuffer, error) {
	byteBuffer := NewMessageByteBuffer()
	byteBuffer.Append(rpcProtocol.MagicBytes)

	rpcMessage, _ := message.(*RpcMessage)
	var data []byte
	dataLength := 0

	switch rpcMessage.DataType {
	case TEXT:
		if d, ok := rpcMessage.Data.(string); ok {
			data = []byte(d)
			dataLength = len(d)
		}
	case GZIP_TEXT:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	case BIN:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	case GZIP_BIN:
		if d, ok := rpcMessage.Data.([]byte); ok {
			data = d
			dataLength = len(d)
		}
	default:
		return nil, errors.New("not support")
	}
	byteBuffer.Append(Int64ToBytes(29 + int64(dataLength) + int64(rpcProtocol.GetCheckCodeLength())))
	byteBuffer.Append(rpcMessage.MessageType)
	byteBuffer.Append(Int64ToBytes(rpcMessage.Id))
	byteBuffer.Append(Int64ToBytes(rpcMessage.ResponseId))
	byteBuffer.AppendByte(rpcMessage.DataType)
	byteBuffer.Append(data)
	if rpcProtocol.CheckCodeBytes != nil {
		byteBuffer.Append(rpcProtocol.CheckCodeBytes)
	}

	return byteBuffer, nil
}

type RpcContext struct {
	RequestId    int64
	RequestTime  time.Time
	ResponseTime time.Time
	Response     chan *RpcMessage
}

func NewRpcContext(requestId int64) *RpcContext {
	return &RpcContext{
		RequestId:    requestId,
		ResponseTime: time.Now(),
		Response:     make(chan *RpcMessage, 1),
	}
}

type RpcTemplate struct {
	session        *Session
	rpcContexts    *sync.Map
	messageHandles map[string]func(message *RpcMessage)
}

func NewRpcTemplate(session *Session) *RpcTemplate {
	return &RpcTemplate{
		rpcContexts:    &sync.Map{},
		messageHandles: make(map[string]func(message *RpcMessage)),
		session:        session,
	}
}

func (rpc *RpcTemplate) OnMessage(message interface{}, session *Session) {
	if rpcMessage, ok := (message).(*RpcMessage); ok {
		if requestContext := rpc.GetRpcContext(rpcMessage.ResponseId); requestContext != nil {
			requestContext.ResponseTime = time.Now()
			requestContext.Response <- rpcMessage
		}
		go func() {
			if f := rpc.messageHandles[hex.EncodeToString(rpcMessage.MessageType)]; f != nil {
				f(rpcMessage)
			}
		}()
	}
}

func (rpc *RpcTemplate) RegisterRpcMessageHandle(messageType string, f func(message *RpcMessage)) {
	rpc.messageHandles[messageType] = f
}

func (rpc *RpcTemplate) Send(m *RpcMessage, timeout time.Duration) (int, error) {
	if m.Id == 0 {
		m.Id = generateRequestIndex()
	}
	return rpc.session.SendMessage(m)
}

var requestIndex int64 = 0

func generateRequestIndex() int64 {
	return atomic.AddInt64(&requestIndex, 1)
}
func (rpc *RpcTemplate) SendWithResponse(m *RpcMessage, timeout time.Duration) (*RpcMessage, error) {
	m.Id = generateRequestIndex()
	rpcContext := NewRpcContext(m.Id)
	rpc.SetRpcContext(m.Id, rpcContext)
	_, err := rpc.Send(m, timeout)
	if err != nil {
		return nil, err
	}
	response := <-rpcContext.Response
	rpc.RemoveRpcContext(m.Id)
	return response, nil
}

func (rpc *RpcTemplate) SetRpcContext(requestId int64, ctx *RpcContext) {
	rpc.rpcContexts.Store(requestId, ctx)
}

func (rpc *RpcTemplate) GetRpcContext(requestId int64) *RpcContext {
	data, _ := rpc.rpcContexts.Load(requestId)
	ctx, _ := data.(*RpcContext)
	return ctx
}

func (rpc *RpcTemplate) RemoveRpcContext(requestId int64) {
	rpc.rpcContexts.Delete(requestId)
}
