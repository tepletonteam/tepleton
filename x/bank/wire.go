package bank

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	// TODO include option to always include prefix bytes.
	//cdc.RegisterConcrete(SendMsg{}, "github.com/tepleton/tepleton-sdk/bank/SendMsg", nil)
	//cdc.RegisterConcrete(IssueMsg{}, "github.com/tepleton/tepleton-sdk/bank/IssueMsg", nil)
}
