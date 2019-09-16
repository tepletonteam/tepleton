package abi

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/errors"
)

// nolint
var (
	errChainNotRegistered  = fmt.Errorf("Chain not registered")
	errChainAlreadyExists  = fmt.Errorf("Chain already exists")
	errNeedsABIPermission  = fmt.Errorf("Needs app-permission to send ABI")
	errCannotSetPermission = fmt.Errorf("Requesting invalid permission on ABI")
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
	ABICodeInvalidCall         = wrsp.CodeType(1007)
)

func ErrNotRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainNotRegistered, ABICodeChainNotRegistered)
}
func IsNotRegisteredErr(err error) bool {
	return errors.IsSameError(errChainNotRegistered, err)
}

func ErrAlreadyRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainAlreadyExists, ABICodeChainAlreadyExists)
}
func IsAlreadyRegistetedErr(err error) bool {
	return errors.IsSameError(errChainAlreadyExists, err)
}

func ErrNeedsABIPermission() error {
	return errors.WithCode(errNeedsABIPermission, ABICodeInvalidCall)
}
func IsNeedsABIPermissionErr(err error) bool {
	return errors.IsSameError(errNeedsABIPermission, err)
}

func ErrCannotSetPermission() error {
	return errors.WithCode(errCannotSetPermission, ABICodeInvalidCall)
}
func IsCannotSetPermissionErr(err error) bool {
	return errors.IsSameError(errCannotSetPermission, err)
}
