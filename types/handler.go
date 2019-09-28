package types

// Handler handles both WRSP DeliverTx and CheckTx requests.
// Iff WRSP.CheckTx, ctx.IsCheckTx() returns true.
type Handler func(ctx Context, ms MultiStore, tx Tx) Result
