package types

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "tepleton-sdk/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "tepleton-sdk/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "tepleton-sdk/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "tepleton-sdk/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgCompleteUnbonding{}, "tepleton-sdk/CompleteUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "tepleton-sdk/BeginRedelegate", nil)
	cdc.RegisterConcrete(MsgCompleteRedelegate{}, "tepleton-sdk/CompleteRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
	//MsgCdc = cdc.Seal() //TODO use when upgraded to go-amino 0.9.10
}
