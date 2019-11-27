package bank

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "tepleton-sdk/Send", nil)
	cdc.RegisterConcrete(MsgIssue{}, "tepleton-sdk/Issue", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
