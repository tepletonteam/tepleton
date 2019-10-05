package types

import wrsp "github.com/tepleton/wrsp/types"

// initialize application state at genesis
type InitChainer func(ctx Context, req wrsp.RequestInitChain) wrsp.ResponseInitChain

// run code before the transactions in a block
type BeginBlocker func(ctx Context, req wrsp.RequestBeginBlock) wrsp.ResponseBeginBlock

// run code after the transactions in a block and return updates to the validator set
type EndBlocker func(ctx Context, req wrsp.RequestEndBlock) wrsp.ResponseEndBlock
