package types

import (
	"fmt"

	cmn "github.com/tepleton/tmlibs/common"

	wrsp "github.com/tepleton/wrsp/types"
)

// WRSPCodeType - combined codetype / codespace
type WRSPCodeType uint32

// CodeType - code identifier within codespace
type CodeType uint16

// CodespaceType - codespace identifier
type CodespaceType uint16

// IsOK - is everything okay?
func (code WRSPCodeType) IsOK() bool {
	if code == WRSPCodeOK {
		return true
	}
	return false
}

func ToWRSPCode(space CodespaceType, code CodeType) WRSPCodeType {
	// TODO: Make Tendermint more aware of codespaces.
	if space == CodespaceRoot && code == CodeOK {
		return WRSPCodeOK
	}
	return WRSPCodeType((uint32(space) << 16) | uint32(code))
}

const (
	// WRSP error codes
	WRSPCodeOK WRSPCodeType = 0

	// Base error codes
	CodeOK                CodeType = 0
	CodeInternal          CodeType = 1
	CodeTxDecode          CodeType = 2
	CodeInvalidSequence   CodeType = 3
	CodeUnauthorized      CodeType = 4
	CodeInsufficientFunds CodeType = 5
	CodeUnknownRequest    CodeType = 6
	CodeInvalidAddress    CodeType = 7
	CodeInvalidPubKey     CodeType = 8
	CodeUnknownAddress    CodeType = 9
	CodeInsufficientCoins CodeType = 10
	CodeInvalidCoins      CodeType = 11

	// CodespaceRoot is a codespace for error codes in this file only.
	// Notice that 0 is an "unset" codespace, which can be overridden with
	// Error.WithDefaultCodespace().
	CodespaceUndefined CodespaceType = 0
	CodespaceRoot      CodespaceType = 1

	// Maximum reservable codespace (2^16 - 1)
	MaximumCodespace CodespaceType = 65535
)

// NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInternal:
		return "Internal error"
	case CodeTxDecode:
		return "Tx parse error"
	case CodeInvalidSequence:
		return "Invalid sequence"
	case CodeUnauthorized:
		return "Unauthorized"
	case CodeInsufficientFunds:
		return "Insufficent funds"
	case CodeUnknownRequest:
		return "Unknown request"
	case CodeInvalidAddress:
		return "Invalid address"
	case CodeInvalidPubKey:
		return "Invalid pubkey"
	case CodeUnknownAddress:
		return "Unknown address"
	case CodeInsufficientCoins:
		return "Insufficient coins"
	case CodeInvalidCoins:
		return "Invalid coins"
	default:
		return fmt.Sprintf("Unknown code %d", code)
	}
}

//--------------------------------------------------------------------------------
// All errors are created via constructors so as to enable us to hijack them
// and inject stack traces if we really want to.

// nolint
func ErrInternal(msg string) Error {
	return newErrorWithRootCodespace(CodeInternal, msg)
}
func ErrTxDecode(msg string) Error {
	return newErrorWithRootCodespace(CodeTxDecode, msg)
}
func ErrInvalidSequence(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidSequence, msg)
}
func ErrUnauthorized(msg string) Error {
	return newErrorWithRootCodespace(CodeUnauthorized, msg)
}
func ErrInsufficientFunds(msg string) Error {
	return newErrorWithRootCodespace(CodeInsufficientFunds, msg)
}
func ErrUnknownRequest(msg string) Error {
	return newErrorWithRootCodespace(CodeUnknownRequest, msg)
}
func ErrInvalidAddress(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidAddress, msg)
}
func ErrUnknownAddress(msg string) Error {
	return newErrorWithRootCodespace(CodeUnknownAddress, msg)
}
func ErrInvalidPubKey(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidPubKey, msg)
}
func ErrInsufficientCoins(msg string) Error {
	return newErrorWithRootCodespace(CodeInsufficientCoins, msg)
}
func ErrInvalidCoins(msg string) Error {
	return newErrorWithRootCodespace(CodeInvalidCoins, msg)
}

//----------------------------------------
// Error & sdkError

// sdk Error type
type Error interface {
	Error() string
	Code() CodeType
	Codespace() CodespaceType
	WRSPLog() string
	WRSPCode() WRSPCodeType
	WithDefaultCodespace(codespace CodespaceType) Error
	Trace(msg string) Error
	T() interface{}
	Result() Result
	QueryResult() wrsp.ResponseQuery
}

// NewError - create an error
func NewError(codespace CodespaceType, code CodeType, msg string) Error {
	return newError(codespace, code, msg)
}

func newErrorWithRootCodespace(code CodeType, msg string) *sdkError {
	return newError(CodespaceRoot, code, msg)
}

func newError(codespace CodespaceType, code CodeType, msg string) *sdkError {
	if msg == "" {
		msg = CodeToDefaultMsg(code)
	}
	return &sdkError{
		codespace: codespace,
		code:      code,
		err:       cmn.NewErrorWithT(code, msg),
	}
}

type sdkError struct {
	codespace CodespaceType
	code      CodeType
	err       cmn.Error
}

// Implements WRSPError.
func (err *sdkError) Error() string {
	return fmt.Sprintf("Error{%d:%d,%#v}", err.codespace, err.code, err.err)
}

// Implements WRSPError.
func (err *sdkError) WRSPCode() WRSPCodeType {
	return ToWRSPCode(err.codespace, err.code)
}

// Implements Error.
func (err *sdkError) Codespace() CodespaceType {
	return err.codespace
}

// Implements Error.
func (err *sdkError) Code() CodeType {
	return err.code
}

// Implements WRSPError.
func (err *sdkError) WRSPLog() string {
	return fmt.Sprintf(`=== WRSP Log ===
Codespace: %v
Code:      %v
WRSPCode:  %v
Error:     %#v
=== /WRSP Log ===
`, err.codespace, err.code, err.WRSPCode(), err.err)
}

// Add tracing information with msg.
func (err *sdkError) Trace(msg string) Error {
	return &sdkError{
		codespace: err.codespace,
		code:      err.code,
		err:       err.err.Trace(msg),
	}
}

// Implements Error.
func (err *sdkError) WithDefaultCodespace(cs CodespaceType) Error {
	codespace := err.codespace
	if codespace == CodespaceUndefined {
		codespace = cs
	}
	return &sdkError{
		codespace: codespace,
		code:      err.code,
		err:       err.err,
	}
}

func (err *sdkError) T() interface{} {
	return err.err.T()
}

func (err *sdkError) Result() Result {
	return Result{
		Code: err.WRSPCode(),
		Log:  err.WRSPLog(),
	}
}

// QueryResult allows us to return sdk.Error.QueryResult() in query responses
func (err *sdkError) QueryResult() wrsp.ResponseQuery {
	return wrsp.ResponseQuery{
		Code: uint32(err.WRSPCode()),
		Log:  err.WRSPLog(),
	}
}
