package simplestake

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgBond{}, "simplestake/BondMsg", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "simplestake/UnbondMsg", nil)
}
