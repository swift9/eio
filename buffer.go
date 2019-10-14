package eio

type ByteBuffer struct {
	buf []byte
}

func (buf *ByteBuffer) Len() int {
	return len(buf.buf)
}

func (buf *ByteBuffer) Discard(len int) {
	buf.buf = buf.buf[len:]
}

func (buf *ByteBuffer) Append(bytes []byte) {
	buf.buf = append(buf.buf, bytes...)
}

func (buf *ByteBuffer) Read(start int, end int) []byte {
	return buf.buf[start:end]
}
