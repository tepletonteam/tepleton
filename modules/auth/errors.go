//nolint
package auth

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	"github.com/tepleton/tepleton-sdk/errors"
)

var (
	errInvalidSignature  = fmt.Errorf("Invalid Signature")   //move auth
	errTooManySignatures = fmt.Errorf("Too many signatures") //move auth

	unauthorized = wrsp.CodeType_Unauthorized
)

func ErrTooManySignatures() errors.TMError {
	return errors.WithCode(errTooManySignatures, unauthorized)
}
func IsTooManySignaturesErr(err error) bool {
	return errors.IsSameError(errTooManySignatures, err)
}

func ErrInvalidSignature() errors.TMError {
	return errors.WithCode(errInvalidSignature, unauthorized)
}
func IsInvalidSignatureErr(err error) bool {
	return errors.IsSameError(errInvalidSignature, err)
}
