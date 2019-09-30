package bank

import (
	"github.com/tepleton/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
	// TODO include option to always include prefix bytes.
	cdc.RegisterConcrete(SendMsg{}, "tepleton-sdk/SendMsg", nil)
	cdc.RegisterConcrete(IssueMsg{}, "tepleton-sdk/IssueMsg", nil)
}
