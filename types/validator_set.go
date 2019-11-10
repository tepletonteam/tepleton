package types

import (
	wrsp "github.com/tepleton/wrsp/types"
)

type Validator = wrsp.Validator

type ValidatorSetKeeper interface {
	Validators(Context) []*Validator
	Size(Context) int
	IsValidator(Context, Address) bool
	GetByAddress(Context, Address) (int, *Validator)
	GetByIndex(Context, int) *Validator
	TotalPower(Context) Rat
}
