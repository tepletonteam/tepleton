package types

import (
	"github.com/tepleton/tepleton-sdk/store"
)

// Handler handles both WRSP DeliverTx and CheckTx requests.
// Iff WRSP.CheckTx, ctx.IsCheckTx() returns true.
type Handler func(ctx Context, store store.MultiStore, tx Tx) Result
