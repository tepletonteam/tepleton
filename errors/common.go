package errors

/**
*    Copyright (C) 2017 Ethan Frey
**/

import (
	rawerr "errors"
	"fmt"

	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin"
)

var (
	errDecoding          = rawerr.New("Error decoding input")
	errUnauthorized      = rawerr.New("Unauthorized")
	errInvalidSignature  = rawerr.New("Invalid Signature")
	errTooLarge          = rawerr.New("Input size too large")
	errMissingSignature  = rawerr.New("Signature missing")
	errTooManySignatures = rawerr.New("Too many signatures")
	errNoChain           = rawerr.New("No chain id provided")
	errWrongChain        = rawerr.New("Wrong chain for tx")
	errUnknownTxType     = rawerr.New("Tx type unknown")
	errInvalidFormat     = rawerr.New("Invalid format")
	errUnknownModule     = rawerr.New("Unknown module")
)

func ErrUnknownTxType(tx basecoin.Tx) TMError {
	msg := fmt.Sprintf("%T", tx.Unwrap())
	w := errors.Wrap(errUnknownTxType, msg)
	return WithCode(w, wrsp.CodeType_UnknownRequest)
}
func IsUnknownTxTypeErr(err error) bool {
	return IsSameError(errUnknownTxType, err)
}

func ErrInvalidFormat(tx basecoin.Tx) TMError {
	msg := fmt.Sprintf("%T", tx.Unwrap())
	w := errors.Wrap(errInvalidFormat, msg)
	return WithCode(w, wrsp.CodeType_UnknownRequest)
}
func IsInvalidFormatErr(err error) bool {
	return IsSameError(errInvalidFormat, err)
}

func ErrUnknownModule(mod string) TMError {
	w := errors.Wrap(errUnknownModule, mod)
	return WithCode(w, wrsp.CodeType_UnknownRequest)
}
func IsUnknownModuleErr(err error) bool {
	return IsSameError(errUnknownModule, err)
}

func ErrInternal(msg string) TMError {
	return New(msg, wrsp.CodeType_InternalError)
}

// IsInternalErr matches any error that is not classified
func IsInternalErr(err error) bool {
	return HasErrorCode(err, wrsp.CodeType_InternalError)
}

func ErrDecoding() TMError {
	return WithCode(errDecoding, wrsp.CodeType_EncodingError)
}
func IsDecodingErr(err error) bool {
	return IsSameError(errDecoding, err)
}

func ErrUnauthorized() TMError {
	return WithCode(errUnauthorized, wrsp.CodeType_Unauthorized)
}

// IsUnauthorizedErr is generic helper for any unauthorized errors,
// also specific sub-types
func IsUnauthorizedErr(err error) bool {
	return HasErrorCode(err, wrsp.CodeType_Unauthorized)
}

func ErrMissingSignature() TMError {
	return WithCode(errMissingSignature, wrsp.CodeType_Unauthorized)
}
func IsMissingSignatureErr(err error) bool {
	return IsSameError(errMissingSignature, err)
}

func ErrTooManySignatures() TMError {
	return WithCode(errTooManySignatures, wrsp.CodeType_Unauthorized)
}
func IsTooManySignaturesErr(err error) bool {
	return IsSameError(errTooManySignatures, err)
}

func ErrInvalidSignature() TMError {
	return WithCode(errInvalidSignature, wrsp.CodeType_Unauthorized)
}
func IsInvalidSignatureErr(err error) bool {
	return IsSameError(errInvalidSignature, err)
}

func ErrNoChain() TMError {
	return WithCode(errNoChain, wrsp.CodeType_Unauthorized)
}
func IsNoChainErr(err error) bool {
	return IsSameError(errNoChain, err)
}

func ErrWrongChain(chain string) TMError {
	msg := errors.Wrap(errWrongChain, chain)
	return WithCode(msg, wrsp.CodeType_Unauthorized)
}
func IsWrongChainErr(err error) bool {
	return IsSameError(errWrongChain, err)
}

func ErrTooLarge() TMError {
	return WithCode(errTooLarge, wrsp.CodeType_EncodingError)
}
func IsTooLargeErr(err error) bool {
	return IsSameError(errTooLarge, err)
}
