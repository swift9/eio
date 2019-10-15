package rpc

type EioMessage struct {
	MessageType          []byte
	MessageSize          int64
	RequestId            int64
	ResponseId           int64
	DataContentType      byte
	CompressionAlgorithm byte
	Body                 interface{}
}
