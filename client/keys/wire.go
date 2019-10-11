package keys

import (
	crypto "github.com/tepleton/go-crypto"
	wire "github.com/tepleton/go-wire"
)

var cdc *wire.Codec

func init() {
	cdc = wire.NewCodec()
	crypto.RegisterWire(cdc)
}

func MarshalJSON(o interface{}) ([]byte, error) {
	return cdc.MarshalJSON(o)
}
