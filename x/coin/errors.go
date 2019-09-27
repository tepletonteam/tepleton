//nolint
package coin

import (
	"fmt"

	"github.com/tepleton/tepleton-sdk/errors"
)

var (
	errNoAccount          = fmt.Errorf("No such account")
	errInsufficientFunds  = fmt.Errorf("Insufficient funds")
	errInsufficientCredit = fmt.Errorf("Insufficient credit")
	errNoInputs           = fmt.Errorf("No input coins")
	errNoOutputs          = fmt.Errorf("No output coins")
	errInvalidAddress     = fmt.Errorf("Invalid address")
	errInvalidCoins       = fmt.Errorf("Invalid coins")
)

const (
	CodeInvalidInput   uint32 = 101
	CodeInvalidOutput  uint32 = 102
	CodeUnknownAddress uint32 = 103
	CodeUnknownRequest uint32 = errors.CodeUnknownRequest
)

func ErrNoAccount() errors.WRSPError {
	return errors.WithCode(errNoAccount, CodeUnknownAddress)
}

func ErrInvalidAddress() errors.WRSPError {
	return errors.WithCode(errInvalidAddress, CodeInvalidInput)
}

func ErrInvalidCoins() errors.WRSPError {
	return errors.WithCode(errInvalidCoins, CodeInvalidInput)
}

func ErrInsufficientFunds() errors.WRSPError {
	return errors.WithCode(errInsufficientFunds, CodeInvalidInput)
}

func ErrInsufficientCredit() errors.WRSPError {
	return errors.WithCode(errInsufficientCredit, CodeInvalidInput)
}

func ErrNoInputs() errors.WRSPError {
	return errors.WithCode(errNoInputs, CodeInvalidInput)
}

func ErrNoOutputs() errors.WRSPError {
	return errors.WithCode(errNoOutputs, CodeInvalidOutput)
}
