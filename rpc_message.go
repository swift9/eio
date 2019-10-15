package eio

type RpcMessage struct {
	MessageSize     int64
	MessageType     []byte
	RequestId       int64
	ResponseId      int64
	DataContentType byte
	Body            interface{}
}
