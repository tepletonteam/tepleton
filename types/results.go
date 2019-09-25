package types

import (
	wrsp "github.com/tepleton/wrsp/types"
)

// CheckResult captures any non-error WRSP  result
// to make sure people use error for error cases.
type CheckResult struct {
	wrsp.Result

	// GasAllocated is the maximum units of work we allow this tx to perform
	GasAllocated uint64

	// GasPayment is the total fees for this tx (or other source of payment)
	GasPayment uint64
}

// DeliverResult captures any non-error wrsp result
// to make sure people use error for error cases
type DeliverResult struct {
	wrsp.Result

	// TODO comment
	Diff []*wrsp.Validator

	// TODO comment
	GasUsed uint64
}
