package stake

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "tepleton-sdk/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "tepleton-sdk/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "tepleton-sdk/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "tepleton-sdk/MsgUnbond", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
