package stack

import (
	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/state"
)

//nolint
const (
	NameCheckpoint = "check"
)

// Checkpoint isolates all data store below this
type Checkpoint struct {
	PassOption
}

// Name of the module - fulfills Middleware interface
func (Checkpoint) Name() string {
	return NameCheckpoint
}

var _ Middleware = Checkpoint{}

// CheckTx reverts all data changes if there was an error
func (c Checkpoint) CheckTx(ctx basecoin.Context, store state.KVStore, tx basecoin.Tx, next basecoin.Checker) (res basecoin.Result, err error) {
	ps := state.NewKVCache(unwrap(store))
	res, err = next.CheckTx(ctx, ps, tx)
	if err == nil {
		ps.Sync()
	}
	return res, err
}

// DeliverTx reverts all data changes if there was an error
func (c Checkpoint) DeliverTx(ctx basecoin.Context, store state.KVStore, tx basecoin.Tx, next basecoin.Deliver) (res basecoin.Result, err error) {
	ps := state.NewKVCache(unwrap(store))
	res, err = next.DeliverTx(ctx, ps, tx)
	if err == nil {
		ps.Sync()
	}
	return res, err
}