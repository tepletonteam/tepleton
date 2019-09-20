package eyes

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	"github.com/tepleton/basecoin/errors"
)

var (
	errMissingData = fmt.Errorf("All tx fields must be filled")

	malformed = wrsp.CodeType_EncodingError
)

//nolint
func ErrMissingData() errors.TMError {
	return errors.WithCode(errMissingData, malformed)
}
func IsMissingDataErr(err error) bool {
	return errors.IsSameError(errMissingData, err)
}
