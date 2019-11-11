package types

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-crypto"
)

type Validator interface {
	GetAddress() Address
	GetPubKey() crypto.PubKey
	GetPower() Rat
}

func WRSPValidator(v Validator) wrsp.Validator {
	return wrsp.Validator{
		PubKey: v.GetPubKey().Bytes(),
		Power:  v.GetPower().Evaluate(),
	}
}

type ValidatorSet interface {
	Iterate(func(int, Validator))
	Size() int
}

type ValidatorSetKeeper interface {
	ValidatorSet(Context) ValidatorSet
	Validator(Context, Address) Validator
	TotalPower(Context) Rat
	DelegationSet(Context, Address) DelegationSet
	Delegation(Context, Address, Address) Delegation
}

type Delegation interface {
	GetDelegator() Address
	GetValidator() Address
	GetBondAmount() Rat
}

type DelegationSet interface {
	Iterate(func(int, Delegation))
	Size() int
}
