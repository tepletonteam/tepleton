package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// WRSP Response Codes
	// Base SDK reserves 0 ~ 99.
	CodeInternalError     uint32 = 1
	CodeTxParseError             = 2
	CodeBadNonce                 = 3
	CodeUnauthorized             = 4
	CodeInsufficientFunds        = 5
	CodeUnknownRequest           = 6
)

// NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultLog(code uint32) string {
	switch code {
	case CodeInternalError:
		return "Internal error"
	case CodeTxParseError:
		return "Tx parse error"
	case CodeBadNonce:
		return "Bad nonce"
	case CodeUnauthorized:
		return "Unauthorized"
	case CodeInsufficientFunds:
		return "Insufficent funds"
	case CodeUnknownRequest:
		return "Unknown request"
	default:
		return fmt.Sprintf("Unknown code %d", code)
	}
}

//--------------------------------------------------------------------------------
// All errors are created via constructors so as to enable us to hijack them
// and inject stack traces if we really want to.

func InternalError(log string) *sdkError {
	return newSDKError(CodeInternalError, log)
}

func TxParseError(log string) *sdkError {
	return newSDKError(CodeTxParseError, log)
}

func BadNonce(log string) *sdkError {
	return newSDKError(CodeBadNonce, log)
}

func Unauthorized(log string) *sdkError {
	return newSDKError(CodeUnauthorized, log)
}

func InsufficientFunds(log string) *sdkError {
	return newSDKError(CodeInsufficientFunds, log)
}

func UnknownRequest(log string) *sdkError {
	return newSDKError(CodeUnknownRequest, log)
}

//----------------------------------------
// WRSPError & sdkError

type WRSPError interface {
	WRSPCode() uint32
	WRSPLog() string
	Error() string
}

func NewWRSPError(code uint32, log string) WRSPError {
	return newSDKError(code, log)
}

/*

	This struct is intended to be used with pkg/errors.

	Usage:

	```
		import sdk "github.com/tepleton/tepleton-sdk"
		import "github.com/pkg/errors"

		var err = <some causal error>
		if err != nil {
			err = sdk.InternalError("").WithCause(err)
			err = errors.Wrap(err, "Captured the stack!")
			return err
		}
	```

	Then, to get the WRSP code/log, use WRSPErrorCause()

*/
type sdkError struct {
	code  uint32
	log   string
	cause error
	// TODO stacktrace, optional.
}

func newSDKError(code uint32, log string) *sdkError {
	// TODO capture stacktrace if ENV is set.
	if log == "" {
		log = CodeToDefaultLog(code)
	}
	return &sdkError{
		code:  code,
		log:   log,
		cause: nil,
	}
}

// Implements WRSPError
func (err *sdkError) Error() string {
	return fmt.Sprintf("SDKError{%d: %s}", err.code, err.log)
}

// Implements WRSPError
func (err *sdkError) WRSPCode() uint32 {
	return err.code
}

// Implements WRSPError
func (err *sdkError) WRSPLog() string {
	return err.log
}

// Implements pkg/errors.causer
func (err *sdkError) Cause() error {
	if err.cause != nil {
		return err.cause
	}
	return err
}

// Creates a cloned *sdkError with specific cause
func (err *sdkError) WithCause(cause error) *sdkError {
	copy := *err
	copy.cause = cause
	return &copy
}

//----------------------------------------

// HasSameCause returns true if both errors
// have the same cause.
func HasSameCause(err1 error, err2 error) bool {
	if err1 != nil || err2 != nil {
		panic("HasSomeCause() requires non-nil arguments")
	}
	return Cause(err1) == Cause(err2)
}

// Like Cause but stops upon finding an WRSPError.
// If no error in the cause chain is an WRSPError,
// returns nil.
func WRSPErrorCause(err error) WRSPError {
	for err != nil {
		wrspErr, ok := err.(WRSPError)
		if ok {
			return wrspErr
		}
		cause, ok := err.(causer)
		if !ok {
			return nil
		}
		errCause := cause.Cause()
		if errCause == nil || errCause == err {
			return err
		}
		err = errCause
	}
	return err
}

// Identitical to pkg/errors.Cause, except handles .Cause()
// returning itself.
// TODO: Merge https://github.com/pkg/errors/issues/89 and
// delete this.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			return err
		}
		errCause := cause.Cause()
		if errCause == nil || errCause == err {
			return err
		}
		err = errCause
	}
	return err
}

type causer interface {
	Cause() error
}
