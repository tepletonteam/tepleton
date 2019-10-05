package types

import wrsp "github.com/tepleton/wrsp/types"

// initialize application state at genesis
type InitChainer func(ctx Context, req wrsp.RequestInitChain) Error
