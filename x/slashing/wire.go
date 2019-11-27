package slashing

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgUnrevoke{}, "tepleton-sdk/MsgUnrevoke", nil)
}

var cdcEmpty = wire.NewCodec()
