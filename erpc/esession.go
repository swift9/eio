package erpc

import (
	"encoding/hex"
	"github.com/swift9/eio"
	"sync"
	"sync/atomic"
	"time"
)

type ESession struct {
	Session        *eio.Session
	Seq            int64
	rpcContexts    *sync.Map
	messageHandles map[string]func(message *EMessage, eSession *ESession)
}

func NewESession(session *eio.Session) *ESession {
	return &ESession{
		rpcContexts:    &sync.Map{},
		Seq:            0,
		messageHandles: make(map[string]func(message *EMessage, eSession *ESession)),
		Session:        session,
	}
}

type rpcContext struct {
	RequestId    int64
	RequestTime  time.Time
	ResponseTime time.Time
	Response     chan *EMessage
}

func newRpcContext(requestId int64) *rpcContext {
	return &rpcContext{
		RequestId:    requestId,
		ResponseTime: time.Now(),
		Response:     make(chan *EMessage, 1),
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
				f(rpcMessage, eSession)
			}
		}()
	}
}

func (eSession *ESession) RegisterMessageHandle(messageType string, f func(message *EMessage, session *ESession)) {
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
	rpcContext := newRpcContext(m.Id)
	eSession.setRpcContext(m.Id, rpcContext)
	_, err := eSession.Send(m, timeout)
	if err != nil {
		return nil, err
	}
	response := <-rpcContext.Response
	eSession.destroyRpcContext(rpcContext)
	return response, nil
}

func (eSession *ESession) Close() error {
	err := eSession.Session.Close()

	eSession.rpcContexts.Range(func(key, value interface{}) bool {
		defer func() {
			if err := recover(); err != nil {
				eSession.Session.Log.Error("close ", err)
			}
		}()
		ctx, _ := value.(*rpcContext)
		close(ctx.Response)
		return false
	})

	return err
}

func (eSession *ESession) setRpcContext(requestId int64, ctx *rpcContext) {
	eSession.rpcContexts.Store(requestId, ctx)
}

func (eSession *ESession) getRpcContext(requestId int64) *rpcContext {
	data, _ := eSession.rpcContexts.Load(requestId)
	ctx, _ := data.(*rpcContext)
	return ctx
}

func (eSession *ESession) destroyRpcContext(c *rpcContext) {
	close(c.Response)
	eSession.rpcContexts.Delete(c.RequestId)
}
