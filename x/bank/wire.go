package bank

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(SendMsg{}, "tepleton-sdk/Send", nil)
	cdc.RegisterConcrete(IssueMsg{}, "tepleton-sdk/Issue", nil)
}
