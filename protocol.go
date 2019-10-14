package eio

// 协议处理
// 报文分包、编解码
type Protocol interface {
	// 报文分包
	Segment(buf *ByteBuffer) []byte

	// 解码
	Decode(bytes []byte) (interface{}, error)

	// 编码
	Encode(d interface{}) ([]byte, error)
}
