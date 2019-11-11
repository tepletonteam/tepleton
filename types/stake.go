package types

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-crypto"
)

// status of a validator
type ValidatorStatus byte

// nolint
const (
	Bonded    ValidatorStatus = 0x00
	Unbonding ValidatorStatus = 0x01
	Unbonded  ValidatorStatus = 0x02
	Revoked   ValidatorStatus = 0x03
)

// validator for a delegated proof of stake system
type Validator interface {
	Status() ValidatorStatus  // status of the validator
	GetOwner() Address        // owner address to receive/return validators coins
	GetPubKey() crypto.PubKey // validation pubkey
	GetPower() Rat            // validation power
	GetBondHeight() int64     // height in which the validator became active
}

// validator which fulfills wrsp validator interface for use in Tendermint
func WRSPValidator(v Validator) wrsp.Validator {
	return wrsp.Validator{
		PubKey: v.GetPubKey().Bytes(),
		Power:  v.GetPower().Evaluate(),
	}
}

// properties for the set of all validators
type ValidatorSet interface {
	IterateValidatorsBonded(func(index int64, validator Validator)) // execute arbitrary logic for each validator
	Validator(Context, Address) Validator                           // get a particular validator by owner address
	TotalPower(Context) Rat                                         // total power of the validator set
}

//_______________________________________________________________________________

// delegation bond for a delegated proof of stake system
type Delegation interface {
	GetDelegator() Address // delegator address for the bond
	GetValidator() Address // validator owner address for the bond
	GetBondShares() Rat    // amount of validator's shares
}

// properties for the set of all delegations for a particular
type DelegationSet interface {

	// execute arbitrary logic for each validator which a delegator has a delegation for
	IterateDelegators(delegator Address, fn func(index int64, delegation Delegation))
}
