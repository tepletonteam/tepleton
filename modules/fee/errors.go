//nolint
package fee

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	"github.com/tepleton/basecoin/errors"
)

var (
	errInsufficientFees = fmt.Errorf("Insufficient fees")
	errWrongFeeDenom    = fmt.Errorf("Required fee denomination")

	invalidInput = wrsp.CodeType_BaseInvalidInput
)

func ErrInsufficientFees() errors.TMError {
	return errors.WithCode(errInsufficientFees, invalidInput)
}
func IsInsufficientFeesErr(err error) bool {
	return errors.IsSameError(errInsufficientFees, err)
}

func ErrWrongFeeDenom(denom string) errors.TMError {
	return errors.WithMessage(denom, errWrongFeeDenom, invalidInput)
}
func IsWrongFeeDenomErr(err error) bool {
	return errors.IsSameError(errWrongFeeDenom, err)
}
