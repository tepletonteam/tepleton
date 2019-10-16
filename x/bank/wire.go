package bank

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	// TODO: bring this back ...
	/*
		// TODO include option to always include prefix bytes.
		cdc.RegisterConcrete(SendMsg{}, "tepleton-sdk/SendMsg", nil)
		cdc.RegisterConcrete(IssueMsg{}, "tepleton-sdk/IssueMsg", nil)
	*/
}
