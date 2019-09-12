//nolint
package errors

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
)

var (
	errDecoding          = fmt.Errorf("Error decoding input")
	errUnauthorized      = fmt.Errorf("Unauthorized")
	errInvalidSignature  = fmt.Errorf("Invalid Signature")
	errTooLarge          = fmt.Errorf("Input size too large")
	errMissingSignature  = fmt.Errorf("Signature missing")
	errTooManySignatures = fmt.Errorf("Too many signatures")
	errNoChain           = fmt.Errorf("No chain id provided")
	errWrongChain        = fmt.Errorf("Wrong chain for tx")
	errUnknownTxType     = fmt.Errorf("Tx type unknown")
	errInvalidFormat     = fmt.Errorf("Invalid format")
	errUnknownModule     = fmt.Errorf("Unknown module")
	errExpired           = fmt.Errorf("Tx expired")
)

// some crazy reflection to unwrap any generated struct.
func unwrap(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	m := v.MethodByName("Unwrap")
	if m.IsValid() {
		out := m.Call(nil)
		if len(out) == 1 {
			return out[0].Interface()
		}
	}
	return i
}

func ErrUnknownTxType(tx interface{}) TMError {
	msg := fmt.Sprintf("%T", unwrap(tx))
	w := errors.Wrap(errUnknownTxType, msg)
	return WithCode(w, wrsp.CodeType_UnknownRequest)
}
func IsUnknownTxTypeErr(err error) bool {
	return IsSameError(errUnknownTxType, err)
}

func ErrInvalidFormat(tx interface{}) TMError {
	msg := fmt.Sprintf("%T", unwrap(tx))
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

func ErrExpired() TMError {
	return WithCode(errExpired, wrsp.CodeType_Unauthorized)
}
func IsExpiredErr(err error) bool {
	return IsSameError(errExpired, err)
}
