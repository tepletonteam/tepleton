package abi

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/errors"
)

// nolint
var (
	errChainNotRegistered = fmt.Errorf("Chain not registered")
	errChainAlreadyExists = fmt.Errorf("Chain already exists")
	// errNotMember        = fmt.Errorf("Not a member")
	// errInsufficientSigs = fmt.Errorf("Not enough signatures")
	// errNoMembers        = fmt.Errorf("No members specified")
	// errTooManyMembers   = fmt.Errorf("Too many members specified")
	// errNotEnoughMembers = fmt.Errorf("Not enough members specified")

	ABICodeChainNotRegistered  = wrsp.CodeType(1001)
	ABICodeChainAlreadyExists  = wrsp.CodeType(1002)
	ABICodePacketAlreadyExists = wrsp.CodeType(1003)
	ABICodeUnknownHeight       = wrsp.CodeType(1004)
	ABICodeInvalidCommit       = wrsp.CodeType(1005)
	ABICodeInvalidProof        = wrsp.CodeType(1006)
)

func ErrNotRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainNotRegistered, ABICodeChainNotRegistered)
}
func IsNotRegistetedErr(err error) bool {
	return errors.IsSameError(errChainNotRegistered, err)
}

func ErrAlreadyRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainAlreadyExists, ABICodeChainAlreadyExists)
}
func IsAlreadyRegistetedErr(err error) bool {
	return errors.IsSameError(errChainAlreadyExists, err)
}
