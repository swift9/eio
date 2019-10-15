package eio

type MessageByteBuffer struct {
	buf []byte
	len int64
}

func (messageByteBuffer *MessageByteBuffer) Len() int64 {
	return messageByteBuffer.len
}

func (messageByteBuffer *MessageByteBuffer) Discard(l int64) {
	messageByteBuffer.buf = messageByteBuffer.buf[l:]
	messageByteBuffer.len -= l
}

func (messageByteBuffer *MessageByteBuffer) Append(bytes []byte) {
	messageByteBuffer.buf = append(messageByteBuffer.buf, bytes...)
	messageByteBuffer.len += int64(len(bytes))
}

func (messageByteBuffer *MessageByteBuffer) AppendByte(b byte) {
	messageByteBuffer.buf = append(messageByteBuffer.buf, b)
	messageByteBuffer.len += 1
}

func (messageByteBuffer *MessageByteBuffer) Peek(start int64, end int64) *MessageByteBuffer {
	return &MessageByteBuffer{
		buf: messageByteBuffer.buf[start:end],
	}
}

func (messageByteBuffer *MessageByteBuffer) Int64Value(start int64, end int64) int64 {
	return BytesToInt64(messageByteBuffer.buf[start:end])
}

func (messageByteBuffer *MessageByteBuffer) Message() []byte {
	return messageByteBuffer.buf
}
