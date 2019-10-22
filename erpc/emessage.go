package erpc

import "github.com/swift9/eio"

type EMessage struct {
	Id          int64
	ResponseId  int64
	MessageSize int64
	MessageType []byte
	DataType    byte
	Data        *eio.MessageBuffer
}
