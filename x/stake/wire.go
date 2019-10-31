package stake

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// TODO complete when go-amino is ported
func RegisterWire(cdc *wire.Codec) {
	// TODO include option to always include prefix bytes.
	//cdc.RegisterConcrete(SendMsg{}, "tepleton-sdk/SendMsg", nil)
	//cdc.RegisterConcrete(IssueMsg{}, "tepleton-sdk/IssueMsg", nil)
}
