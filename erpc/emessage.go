package erpc

type EMessage struct {
	Id          int64
	ResponseId  int64
	MessageSize int64
	MessageType []byte
	DataType    byte
	Data        interface{}
}
