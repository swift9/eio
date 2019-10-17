package eio

// 协议处理
// 报文分包、编解码
type Protocol interface {
	// 报文分包
	Segment(session *Session, buf *MessageBuffer) *MessageBuffer

	// 解码
	Decode(session *Session, message *MessageBuffer) (interface{}, error)

	// 编码
	Encode(session *Session, d interface{}) (*MessageBuffer, error)

	// 验证报文有效性
	IsValidMessage(session *Session, message *MessageBuffer) bool
}
