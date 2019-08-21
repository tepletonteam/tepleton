package types

import (
	wrsp "github.com/tepleton/wrsp/types"
)

var (
	ErrInternalError        = wrsp.NewError(wrsp.CodeType_InternalError, "Internal error")
	ErrDuplicateAddress     = wrsp.NewError(wrsp.CodeType_BaseDuplicateAddress, "Error duplicate address")
	ErrEncodingError        = wrsp.NewError(wrsp.CodeType_BaseEncodingError, "Error encoding error")
	ErrInsufficientFees     = wrsp.NewError(wrsp.CodeType_BaseInsufficientFees, "Error insufficient fees")
	ErrInsufficientFunds    = wrsp.NewError(wrsp.CodeType_BaseInsufficientFunds, "Error insufficient funds")
	ErrInsufficientGasPrice = wrsp.NewError(wrsp.CodeType_BaseInsufficientGasPrice, "Error insufficient gas price")
	ErrInvalidAddress       = wrsp.NewError(wrsp.CodeType_BaseInvalidAddress, "Error invalid address")
	ErrInvalidAmount        = wrsp.NewError(wrsp.CodeType_BaseInvalidAmount, "Error invalid amount")
	ErrInvalidPubKey        = wrsp.NewError(wrsp.CodeType_BaseInvalidPubKey, "Error invalid pubkey")
	ErrInvalidSequence      = wrsp.NewError(wrsp.CodeType_BaseInvalidSequence, "Error invalid sequence")
	ErrInvalidSignature     = wrsp.NewError(wrsp.CodeType_BaseInvalidSignature, "Error invalid signature")
	ErrUnknownPubKey        = wrsp.NewError(wrsp.CodeType_BaseUnknownPubKey, "Error unknown pubkey")

	ResultOK = wrsp.NewResultOK(nil, "")
)
