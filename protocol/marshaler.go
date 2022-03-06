package protocol

import (
	"github.com/quantosnetwork/Quantos/decoder"
	"github.com/quantosnetwork/Quantos/encoder"
)

func MarshalTo(dst []byte, data interface{}) ([]byte, error) {
	var e encoder.Encoder
	return e.EncodeTo(dst, data)
}

func Marshal(data interface{}) ([]byte, error) {
	var e encoder.Encoder
	return e.EncodeTo(nil, data)
}

func Unmashal(data []byte) (interface{}, error) {
	var d decoder.Decoder
	return d.Decode(data)
}
