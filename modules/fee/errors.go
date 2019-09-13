//nolint
package fee

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	"github.com/tepleton/basecoin/errors"
)

var (
	errInsufficientFees = fmt.Errorf("Insufficient Fees")
	errWrongFeeDenom    = fmt.Errorf("Required fee denomination")
)

func ErrInsufficientFees() errors.TMError {
	return errors.WithCode(errInsufficientFees, wrsp.CodeType_BaseInvalidInput)
}
func IsInsufficientFeesErr(err error) bool {
	return errors.IsSameError(errInsufficientFees, err)
}

func ErrWrongFeeDenom(denom string) errors.TMError {
	return errors.WithMessage(denom, errWrongFeeDenom, wrsp.CodeType_BaseInvalidInput)
}
func IsWrongFeeDenomErr(err error) bool {
	return errors.IsSameError(errWrongFeeDenom, err)
}
