package rpc

import "github.com/swift9/eio"

type EioProtocol struct {
	eio.VariableProtocol
	CheckCodeBytes []byte
}

func (p *EioProtocol) IsValidMessage(bytes []byte) bool {
	if p.CheckCodeBytes == nil {
		return true
	}
	return true
}

func (p *EioProtocol) Decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (p *EioProtocol) Encode(message interface{}) ([]byte, error) {
	bytes, _ := message.([]byte)
	return bytes, nil
}

func NewRpcProtocol() {

}
