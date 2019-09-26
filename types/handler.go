package types

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/log"
)

// Handler is something that processes a transaction.
type Handler interface {

	// Checker verifies there are valid fees and estimates work.
	CheckTx(ctx Context, ms MultiStore, tx Tx) CheckResult

	// Deliverer performs the tx once it makes it in the block.
	DeliverTx(ctx Context, ms MultiStore, tx Tx) DeliverResult
}

// Checker verifies there are valid fees and estimates work.
// NOTE: Keep in sync with Handler.CheckTx
type CheckTxFunc func(ctx Context, ms MultiStore, tx Tx) CheckResult

// Deliverer performs the tx once it makes it in the block.
// NOTE: Keep in sync with Handler.DeliverTx
type DeliverTxFunc func(ctx Context, ms MultiStore, tx Tx) DeliverResult
