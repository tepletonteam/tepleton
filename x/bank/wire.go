package bank

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	crypto "github.com/tepleton/go-crypto"
	"github.com/tepleton/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
	// TODO include option to always include prefix bytes.
	cdc.RegisterConcrete(SendMsg{}, "tepleton-sdk/SendMsg", nil)
	cdc.RegisterConcrete(IssueMsg{}, "tepleton-sdk/IssueMsg", nil)

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature] types.
}
