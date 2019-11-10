package types

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-crypto"

	"github.com/tepleton/tepleton-sdk/wire"
)

var cdc = wire.NewCodec()

func init() {
	crypto.RegisterAmino(cdc)
}

type Validator = wrsp.Validator

type ValidatorSetKeeper interface {
	GetValidators(Context) []*Validator
	Size(Context) int
	IsValidator(Context, Address) bool
	GetByAddress(Context, Address) (int, *Validator)
	GetByIndex(Context, int) *Validator
	TotalPower(Context) Rat
}
