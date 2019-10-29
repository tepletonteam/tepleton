package server

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

var cdc *wire.Codec

func init() {
	cdc = wire.NewCodec()
	wire.RegisterCrypto(cdc)
}
