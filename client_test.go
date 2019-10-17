package eio_test

import (
	"github.com/swift9/eio"
	"sync"
	"testing"
	"time"
)

func test(threadCount int64, msgCountPerThread int64) time.Time {
	w := &sync.WaitGroup{}
	protocol := &eio.RpcProtocol{}
	protocol.MagicBytes = []byte{0xA0, 0xA0}
	protocol.MessageByteSize = 8
	protocol.CheckCodeBytes = []byte{0x0A, 0x0A}

	start := time.Now()
	client := eio.NewClient("localhost:8000", protocol)
	var rpc *eio.RpcTemplate
	client.Connect(func(session *eio.Session) {
		rpc = eio.NewRpcTemplate(session)
		session.OnMessage = rpc.OnMessage
	})

	w.Add(eio.Int642Int(threadCount))
	for i := 0; i < eio.Int642Int(threadCount); i++ {
		go func() {
			var i int64 = 0
			for {
				i++
				t := time.Now()

				m, _ := rpc.SendWithResponse(&eio.RpcMessage{
					MessageType: []byte{0x00, 0x01},
					DataType:    eio.TEXT,
					Data:        "hello",
				}, 1*time.Second)

				if m.ResponseId%10000 == 0 {
					println(time.Now().String(), m.Id, (time.Now().UnixNano()-t.UnixNano())/1000.0)
				}

				if i == msgCountPerThread+1 {
					break
				}
			}
			w.Done()
		}()
	}
	w.Wait()
	return start
}

func TestClient_Rpc(t *testing.T) {
	println(time.Now().String())
	var threadCount, msgCountPerThread int64 = 16, 20000
	start := test(threadCount, msgCountPerThread)
	println(time.Now().String(), " qps:", (threadCount*msgCountPerThread)/(time.Now().Unix()-start.Unix()))

}
