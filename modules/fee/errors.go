//nolint
package fee

import (
	rawerr "errors"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/errors"
)

var (
	errInsufficientFees = rawerr.New("Insufficient Fees")
)

func ErrInsufficientFees() errors.TMError {
	return errors.WithCode(errInsufficientFees, wrsp.CodeType_BaseInvalidInput)
}
func IsInsufficientFeesErr(err error) bool {
	return errors.IsSameError(errInsufficientFees, err)
}
