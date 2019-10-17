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

func newESessionContext(requestId int64) *ESessionContext {
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

func (eSession *ESession) onMessage(message interface{}, session *eio.Session) {
	if rpcMessage, ok := (message).(*EMessage); ok {
		if requestContext := eSession.getRpcContext(rpcMessage.ResponseId); requestContext != nil {
			requestContext.ResponseTime = time.Now()
			requestContext.Response <- rpcMessage
		}
		go func() {
			if f := eSession.messageHandles[hex.EncodeToString(rpcMessage.MessageType)]; f != nil {
				f(rpcMessage)
			}
		}()
	}
}

func (eSession *ESession) RegisterMessageHandle(messageType string, f func(message *EMessage)) {
	eSession.messageHandles[messageType] = f
}

func (eSession *ESession) Send(m *EMessage, timeout time.Duration) (int, error) {
	if m.Id == 0 {
		m.Id = atomic.AddInt64(&eSession.Seq, 1)
	}
	return eSession.Session.SendMessage(m)
}

func (eSession *ESession) SendWithResponse(m *EMessage, timeout time.Duration) (*EMessage, error) {
	m.Id = atomic.AddInt64(&eSession.Seq, 1)
	rpcContext := newESessionContext(m.Id)
	eSession.setRpcContext(m.Id, rpcContext)
	_, err := eSession.Send(m, timeout)
	if err != nil {
		return nil, err
	}
	response := <-rpcContext.Response
	eSession.removeRpcContext(m.Id)
	return response, nil
}

func (eSession *ESession) Close() error {
	return eSession.Session.Close()
}

func (eSession *ESession) setRpcContext(requestId int64, ctx *ESessionContext) {
	eSession.rpcContexts.Store(requestId, ctx)
}

func (eSession *ESession) getRpcContext(requestId int64) *ESessionContext {
	data, _ := eSession.rpcContexts.Load(requestId)
	ctx, _ := data.(*ESessionContext)
	return ctx
}

func (eSession *ESession) removeRpcContext(requestId int64) {
	eSession.rpcContexts.Delete(requestId)
}
