//nolint
package base

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/errors"
)

var (
	errNoChain    = fmt.Errorf("No chain id provided") //move base
	errWrongChain = fmt.Errorf("Wrong chain for tx")   //move base
	errExpired    = fmt.Errorf("Tx expired")           //move base

	unauthorized = wrsp.CodeType_Unauthorized
)

func ErrNoChain() errors.TMError {
	return errors.WithCode(errNoChain, unauthorized)
}
func IsNoChainErr(err error) bool {
	return errors.IsSameError(errNoChain, err)
}
func ErrWrongChain(chain string) errors.TMError {
	msg := pkgerrors.Wrap(errWrongChain, chain)
	return errors.WithCode(msg, unauthorized)
}
func IsWrongChainErr(err error) bool {
	return errors.IsSameError(errWrongChain, err)
}

func ErrExpired() errors.TMError {
	return errors.WithCode(errExpired, unauthorized)
}
func IsExpiredErr(err error) bool {
	return errors.IsSameError(errExpired, err)
}
