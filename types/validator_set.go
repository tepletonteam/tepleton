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
	GetByAddress(Context, Address) Validator
	TotalPower(Context) Rat
}
