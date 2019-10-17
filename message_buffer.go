package eio

type MessageBuffer struct {
	b   []byte
	len int
}

func (buf *MessageBuffer) Len() int {
	return buf.len
}

func (buf *MessageBuffer) Discard(l int) {
	buf.b = buf.b[l:]
	buf.len -= l
}

func (buf *MessageBuffer) Append(b []byte) {
	buf.b = append(buf.b, b...)
	buf.len += len(b)
}

func (buf *MessageBuffer) AppendByte(b byte) {
	buf.Append([]byte{b})
}

func (buf *MessageBuffer) Peek(start int, end int) *MessageBuffer {
	return &MessageBuffer{
		b:   buf.b[start:end],
		len: end - start,
	}
}

// >= start < end
func (buf *MessageBuffer) Int64Value(start int, end int) int64 {
	b := buf.b[start:end]
	return BytesToInt64(b)
}

func (buf *MessageBuffer) Bytes() []byte {
	return buf.b
}

func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		b:   []byte{},
		len: 0,
	}
}

func NewMessageBufferFrom(bs []byte) *MessageBuffer {
	return &MessageBuffer{
		b:   bs,
		len: len(bs),
	}
}
