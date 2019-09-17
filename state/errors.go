//nolint
package state

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/errors"
)

var (
	errNotASubTransaction = fmt.Errorf("Not a sub-transaction")
)

func ErrNotASubTransaction() errors.TMError {
	return errors.WithCode(errNotASubTransaction, wrsp.CodeType_InternalError)
}
func IsNotASubTransactionErr(err error) bool {
	return errors.IsSameError(errNotASubTransaction, err)
}
