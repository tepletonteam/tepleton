package ibc

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	//cdc.RegisterConcrete(IBCTransferMsg{}, "github.com/tepleton/tepleton-sdk/x/ibc/IBCTransferMsg", nil)
	//cdc.RegisterConcrete(IBCReceiveMsg{}, "github.com/tepleton/tepleton-sdk/x/ibc/IBCReceiveMsg", nil)
}
