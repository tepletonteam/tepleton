package ibc

import (
	wire "github.com/tepleton/go-amino"
)

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "tepleton-sdk/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "tepleton-sdk/IBCReceiveMsg", nil)
	cdc.RegisterConcrete(IBCPacket{}, "tepleton-sdk/IBCPacket", nil)
}
