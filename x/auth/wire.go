package auth

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec for default AppAccount
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "auth/Account", nil)
}
