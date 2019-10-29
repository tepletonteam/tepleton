package rpc

import (
	"github.com/tepleton/tepleton-sdk/wire"
	ctypes "github.com/tepleton/tepleton/rpc/core/types"
)

var cdc *wire.Codec

func init() {
	cdc = wire.NewCodec()
	RegisterWire(cdc)
}

func RegisterWire(cdc *wire.Codec) {
	ctypes.RegisterAmino(cdc)
}
