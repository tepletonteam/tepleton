package types

import wrsp "github.com/tepleton/wrsp/types"

// initialize application state at genesis
type InitChainer func(ctx Context, req wrsp.RequestInitChain) wrsp.ResponseInitChain

//
type BeginBlocker func(ctx Context, req wrsp.RequestBeginBlock) wrsp.ResponseBeginBlock

//
type EndBlocker func(ctx Context, req wrsp.RequestEndBlock) wrsp.ResponseEndBlock
