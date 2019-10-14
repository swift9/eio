package eio

type ByteBuffer struct {
	buf []byte
}

func (buf *ByteBuffer) Len() int64 {
	return int64(len(buf.buf))
}

func (buf *ByteBuffer) Discard(l int) {
	buf.buf = buf.buf[l:]
}

func (buf *ByteBuffer) Append(bytes []byte) {
	buf.buf = append(buf.buf, bytes...)
}

func (buf *ByteBuffer) Read(start int, end int) []byte {
	return buf.buf[start:end]
}
