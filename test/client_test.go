package eio_test

import (
	"github.com/swift9/eio"
	"github.com/swift9/eio/erpc"
	"sync"
	"testing"
	"time"
)

func test(threadCount int64, msgCountPerThread int64) time.Time {
	w := &sync.WaitGroup{}

	start := time.Now()
	client := erpc.NewEClient("localhost:8000", erpc.NewDefaultEProtocol())
	var rpc *erpc.ESession
	client.Connect(func(session *erpc.ESession) {
		rpc = session
	})

	w.Add(eio.Int642Int(threadCount))
	for i := 0; i < eio.Int642Int(threadCount); i++ {
		go func() {
			var i int64 = 0
			for {
				i++
				t := time.Now()

				m, _ := rpc.SendWithResponse(&erpc.EMessage{
					MessageType: []byte{0x00, 0x00, 0x00, 0x01},
					DataType:    erpc.Text,
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
	rpc.Close()
	return start
}

func TestClient_Rpc(t *testing.T) {
	println(time.Now().String())
	var threadCount, msgCountPerThread int64 = 16, 20000
	start := test(threadCount, msgCountPerThread)
	println(time.Now().String(), " qps:", (threadCount*msgCountPerThread)/(time.Now().Unix()-start.Unix()))
}
