package ibc

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
	errNeedsIBCPermission  = fmt.Errorf("Needs app-permission to send IBC")
	errCannotSetPermission = fmt.Errorf("Requesting invalid permission on IBC")
	errHeaderNotFound      = fmt.Errorf("Header not found")
	errPacketAlreadyExists = fmt.Errorf("Packet already handled")
	errPacketOutOfOrder    = fmt.Errorf("Packet out of order")
	errInvalidProof        = fmt.Errorf("Invalid merkle proof")
	msgInvalidCommit       = "Invalid header and commit"

	IBCCodeChainNotRegistered    = wrsp.CodeType(1001)
	IBCCodeChainAlreadyExists    = wrsp.CodeType(1002)
	IBCCodeUnknownChain          = wrsp.CodeType(1003)
	IBCCodeInvalidPacketSequence = wrsp.CodeType(1004)
	IBCCodeUnknownHeight         = wrsp.CodeType(1005)
	IBCCodeInvalidCommit         = wrsp.CodeType(1006)
	IBCCodeInvalidProof          = wrsp.CodeType(1007)
	IBCCodeInvalidCall           = wrsp.CodeType(1008)
)

func ErrNotRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainNotRegistered, IBCCodeChainNotRegistered)
}
func IsNotRegisteredErr(err error) bool {
	return errors.IsSameError(errChainNotRegistered, err)
}

func ErrAlreadyRegistered(chainID string) error {
	return errors.WithMessage(chainID, errChainAlreadyExists, IBCCodeChainAlreadyExists)
}
func IsAlreadyRegisteredErr(err error) bool {
	return errors.IsSameError(errChainAlreadyExists, err)
}

func ErrWrongDestChain(chainID string) error {
	return errors.WithMessage(chainID, errWrongDestChain, IBCCodeUnknownChain)
}
func IsWrongDestChainErr(err error) bool {
	return errors.IsSameError(errWrongDestChain, err)
}

func ErrNeedsIBCPermission() error {
	return errors.WithCode(errNeedsIBCPermission, IBCCodeInvalidCall)
}
func IsNeedsIBCPermissionErr(err error) bool {
	return errors.IsSameError(errNeedsIBCPermission, err)
}

func ErrCannotSetPermission() error {
	return errors.WithCode(errCannotSetPermission, IBCCodeInvalidCall)
}
func IsCannotSetPermissionErr(err error) bool {
	return errors.IsSameError(errCannotSetPermission, err)
}

func ErrHeaderNotFound(h int) error {
	msg := fmt.Sprintf("height %d", h)
	return errors.WithMessage(msg, errHeaderNotFound, IBCCodeUnknownHeight)
}
func IsHeaderNotFoundErr(err error) bool {
	return errors.IsSameError(errHeaderNotFound, err)
}

func ErrPacketAlreadyExists() error {
	return errors.WithCode(errPacketAlreadyExists, IBCCodeInvalidPacketSequence)
}
func IsPacketAlreadyExistsErr(err error) bool {
	return errors.IsSameError(errPacketAlreadyExists, err)
}

func ErrPacketOutOfOrder(seq uint64) error {
	msg := fmt.Sprintf("expected %d", seq)
	return errors.WithMessage(msg, errPacketOutOfOrder, IBCCodeInvalidPacketSequence)
}
func IsPacketOutOfOrderErr(err error) bool {
	return errors.IsSameError(errPacketOutOfOrder, err)
}

func ErrInvalidProof() error {
	return errors.WithCode(errInvalidProof, IBCCodeInvalidProof)
}
func ErrInvalidProofWithReason(err error) error {
	return errors.WithCode(err, IBCCodeInvalidProof)
}
func IsInvalidProofErr(err error) bool {
	return errors.HasErrorCode(err, IBCCodeInvalidProof)
}

func ErrInvalidCommit(err error) error {
	if err == nil {
		return nil
	}
	return errors.WithMessage(msgInvalidCommit, err, IBCCodeInvalidCommit)
}
func IsInvalidCommitErr(err error) bool {
	return errors.HasErrorCode(err, IBCCodeInvalidCommit)
}
