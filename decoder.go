package eio

type IDecoder interface {
	Decode(bytes []byte, dst interface{}) error
}
