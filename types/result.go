package types

import (
	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
)

// Result is the union of ResponseDeliverTx and ResponseCheckTx.
type Result struct {

	// Code is the response code, is stored back on the chain.
	Code WRSPCodeType

	// Data is any data returned from the app.
	Data []byte

	// Log is just debug information. NOTE: nondeterministic.
	Log string

	// GasWanted is the maximum units of work we allow this tx to perform.
	GasWanted int64

	// GasUsed is the amount of gas actually consumed. NOTE: unimplemented
	GasUsed int64

	// Tx fee amount and denom.
	FeeAmount int64
	FeeDenom  string

	// Changes to the validator set.
	ValidatorUpdates []wrsp.Validator

	// Tags are used for transaction indexing and pubsub.
	Tags []cmn.KVPair
}

// TODO: In the future, more codes may be OK.
func (res Result) IsOK() bool {
	return res.Code.IsOK()
}
