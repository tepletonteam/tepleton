package basecoin

import (
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-wire/data"

	"github.com/tepleton/basecoin/types"
)

type Named interface {
	Name() string
}

type Checker interface {
	CheckTx(ctx Context, store types.KVStore, tx Tx) (Result, error)
}

// CheckerFunc (like http.HandlerFunc) is a shortcut for making wrapers
type CheckerFunc func(Context, types.KVStore, Tx) (Result, error)

func (c CheckerFunc) CheckTx(ctx Context, store types.KVStore, tx Tx) (Result, error) {
	return c(ctx, store, tx)
}

type Deliver interface {
	DeliverTx(ctx Context, store types.KVStore, tx Tx) (Result, error)
}

// DeliverFunc (like http.HandlerFunc) is a shortcut for making wrapers
type DeliverFunc func(Context, types.KVStore, Tx) (Result, error)

func (c DeliverFunc) DeliverTx(ctx Context, store types.KVStore, tx Tx) (Result, error) {
	return c(ctx, store, tx)
}

// Handler is anything that processes a transaction
type Handler interface {
	Checker
	Deliver
	Named
	// TODO: flesh these out as well
	// SetOption(store types.KVStore, key, value string) (log string)
	// InitChain(store types.KVStore, vals []*wrsp.Validator)
	// BeginBlock(store types.KVStore, hash []byte, header *wrsp.Header)
	// EndBlock(store types.KVStore, height uint64) wrsp.ResponseEndBlock
}

// Result captures any non-error wrsp result
// to make sure people use error for error cases
type Result struct {
	Data data.Bytes
	Log  string
}

func (r Result) ToWRSP() wrsp.Result {
	return wrsp.Result{
		Data: r.Data,
		Log:  r.Log,
	}
}
