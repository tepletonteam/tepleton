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
	errWrongDestChain      = fmt.Errorf("This is not the destination")
	errNeedsABIPermission  = fmt.Errorf("Needs app-permission to send ABI")
	errCannotSetPermission = fmt.Errorf("Requesting invalid permission on ABI")
	errHeaderNotFound      = fmt.Errorf("Header not found")
	errPacketAlreadyExists = fmt.Errorf("Packet already handled")
	errPacketOutOfOrder    = fmt.Errorf("Packet out of order")
	errInvalidProof        = fmt.Errorf("Invalid merkle proof")
	msgInvalidCommit       = "Invalid header and commit"

	ABICodeChainNotRegistered    = wrsp.CodeType(1001)
	ABICodeChainAlreadyExists    = wrsp.CodeType(1002)
	ABICodeUnknownChain          = wrsp.CodeType(1003)
	ABICodeInvalidPacketSequence = wrsp.CodeType(1004)
	ABICodeUnknownHeight         = wrsp.CodeType(1005)
	ABICodeInvalidCommit         = wrsp.CodeType(1006)
	ABICodeInvalidProof          = wrsp.CodeType(1007)
	ABICodeInvalidCall           = wrsp.CodeType(1008)
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
func IsAlreadyRegisteredErr(err error) bool {
	return errors.IsSameError(errChainAlreadyExists, err)
}

func ErrWrongDestChain(chainID string) error {
	return errors.WithMessage(chainID, errWrongDestChain, ABICodeUnknownChain)
}
func IsWrongDestChainErr(err error) bool {
	return errors.IsSameError(errWrongDestChain, err)
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

func ErrHeaderNotFound(h int) error {
	msg := fmt.Sprintf("height %d", h)
	return errors.WithMessage(msg, errHeaderNotFound, ABICodeUnknownHeight)
}
func IsHeaderNotFoundErr(err error) bool {
	return errors.IsSameError(errHeaderNotFound, err)
}

func ErrPacketAlreadyExists() error {
	return errors.WithCode(errPacketAlreadyExists, ABICodeInvalidPacketSequence)
}
func IsPacketAlreadyExistsErr(err error) bool {
	return errors.IsSameError(errPacketAlreadyExists, err)
}

func ErrPacketOutOfOrder(seq uint64) error {
	msg := fmt.Sprintf("expected %d", seq)
	return errors.WithMessage(msg, errPacketOutOfOrder, ABICodeInvalidPacketSequence)
}
func IsPacketOutOfOrderErr(err error) bool {
	return errors.IsSameError(errPacketOutOfOrder, err)
}

func ErrInvalidProof() error {
	return errors.WithCode(errInvalidProof, ABICodeInvalidProof)
}
func IsInvalidProofErr(err error) bool {
	return errors.IsSameError(errInvalidProof, err)
}

func ErrInvalidCommit(err error) error {
	if err == nil {
		return nil
	}
	return errors.WithMessage(msgInvalidCommit, err, ABICodeInvalidCommit)
}
func IsInvalidCommitErr(err error) bool {
	return errors.HasErrorCode(err, ABICodeInvalidCommit)
}
