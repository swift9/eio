package eio

type MessageByteBuffer struct {
	b   []byte
	len int64
}

func (messageByteBuffer *MessageByteBuffer) Len() int64 {
	return messageByteBuffer.len
}

func (messageByteBuffer *MessageByteBuffer) Discard(l int64) {
	messageByteBuffer.b = messageByteBuffer.b[l:]
	messageByteBuffer.len -= l
}

func (messageByteBuffer *MessageByteBuffer) Append(b []byte) {
	messageByteBuffer.b = append(messageByteBuffer.b, b...)
	messageByteBuffer.len += int64(len(b))
}

func (messageByteBuffer *MessageByteBuffer) AppendByte(b byte) {
	messageByteBuffer.Append([]byte{b})
}

func (messageByteBuffer *MessageByteBuffer) Peek(start int64, end int64) *MessageByteBuffer {
	return &MessageByteBuffer{
		b:   messageByteBuffer.b[start:end],
		len: end - start,
	}
}

func (messageByteBuffer *MessageByteBuffer) Int64Value(start int64, end int64) int64 {
	b := messageByteBuffer.b[start:end]
	return BytesToInt64(b)
}

func (messageByteBuffer *MessageByteBuffer) Message() []byte {
	return messageByteBuffer.b
}

func NewMessageByteBuffer() *MessageByteBuffer {
	return &MessageByteBuffer{
		b:   []byte{},
		len: 0,
	}
}
