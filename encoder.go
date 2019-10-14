package eio

type IEncoder interface {
	Encode(data interface{}) ([]byte, error)
}
