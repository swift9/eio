package erpc

import (
	"encoding/hex"
	"github.com/swift9/eio"
	"sync"
	"sync/atomic"
	"time"
)

type ESessionContext struct {
	RequestId    int64
	RequestTime  time.Time
	ResponseTime time.Time
	Response     chan *EMessage
}

func NewESessionContext(requestId int64) *ESessionContext {
	return &ESessionContext{
		RequestId:    requestId,
		ResponseTime: time.Now(),
		Response:     make(chan *EMessage, 1),
	}
}

type ESession struct {
	Session        *eio.Session
	Seq            int64
	rpcContexts    *sync.Map
	messageHandles map[string]func(message *EMessage)
}

func NewESession(session *eio.Session) *ESession {
	return &ESession{
		rpcContexts:    &sync.Map{},
		Seq:            0,
		messageHandles: make(map[string]func(message *EMessage)),
		Session:        session,
	}
}

func (rpc *ESession) onMessage(message interface{}, session *eio.Session) {
	if rpcMessage, ok := (message).(*EMessage); ok {
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

func (rpc *ESession) RegisterMessageHandle(messageType string, f func(message *EMessage)) {
	rpc.messageHandles[messageType] = f
}

func (rpc *ESession) Send(m *EMessage, timeout time.Duration) (int, error) {
	if m.Id == 0 {
		m.Id = atomic.AddInt64(&rpc.Seq, 1)
	}
	return rpc.Session.SendMessage(m)
}

func (rpc *ESession) SendWithResponse(m *EMessage, timeout time.Duration) (*EMessage, error) {
	m.Id = atomic.AddInt64(&rpc.Seq, 1)
	rpcContext := NewESessionContext(m.Id)
	rpc.SetRpcContext(m.Id, rpcContext)
	_, err := rpc.Send(m, timeout)
	if err != nil {
		return nil, err
	}
	response := <-rpcContext.Response
	rpc.RemoveRpcContext(m.Id)
	return response, nil
}

func (rpc *ESession) SetRpcContext(requestId int64, ctx *ESessionContext) {
	rpc.rpcContexts.Store(requestId, ctx)
}

func (rpc *ESession) GetRpcContext(requestId int64) *ESessionContext {
	data, _ := rpc.rpcContexts.Load(requestId)
	ctx, _ := data.(*ESessionContext)
	return ctx
}

func (rpc *ESession) RemoveRpcContext(requestId int64) {
	rpc.rpcContexts.Delete(requestId)
}
