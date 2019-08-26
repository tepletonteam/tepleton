package counter

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin/types"
	"github.com/tepleton/go-wire"
)

type CounterPluginState struct {
	Counter   int
	TotalFees types.Coins
}

type CounterTx struct {
	Valid bool
	Fee   types.Coins
}

//--------------------------------------------------------------------------------

type CounterPlugin struct {
	name string
}

func (cp *CounterPlugin) Name() string {
	return cp.name
}

func (cp *CounterPlugin) StateKey() []byte {
	return []byte(fmt.Sprintf("CounterPlugin{name=%v}.State", cp.name))
}

func New() *CounterPlugin {
	return &CounterPlugin{
		name: "counter",
	}
}

func (cp *CounterPlugin) SetOption(store types.KVStore, key string, value string) (log string) {
	return ""
}

func (cp *CounterPlugin) RunTx(store types.KVStore, ctx types.CallContext, txBytes []byte) (res wrsp.Result) {
	// Decode tx
	var tx CounterTx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return wrsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error()).PrependLog("CounterTx Error: ")
	}

	// Validate tx
	if !tx.Valid {
		return wrsp.ErrInternalError.AppendLog("CounterTx.Valid must be true")
	}
	if !tx.Fee.IsValid() {
		return wrsp.ErrInternalError.AppendLog("CounterTx.Fee is not sorted or has zero amounts")
	}
	if !tx.Fee.IsNonnegative() {
		return wrsp.ErrInternalError.AppendLog("CounterTx.Fee must be nonnegative")
	}

	// Did the caller provide enough coins?
	if !ctx.Coins.IsGTE(tx.Fee) {
		return wrsp.ErrInsufficientFunds.AppendLog("CounterTx.Fee was not provided")
	}

	// TODO If there are any funds left over, return funds.
	// e.g. !ctx.Coins.Minus(tx.Fee).IsZero()
	// ctx.CallerAccount is synced w/ store, so just modify that and store it.

	// Load CounterPluginState
	var cpState CounterPluginState
	cpStateBytes := store.Get(cp.StateKey())
	if len(cpStateBytes) > 0 {
		err = wire.ReadBinaryBytes(cpStateBytes, &cpState)
		if err != nil {
			return wrsp.ErrInternalError.AppendLog("Error decoding state: " + err.Error())
		}
	}

	// Update CounterPluginState
	cpState.Counter += 1
	cpState.TotalFees = cpState.TotalFees.Plus(tx.Fee)

	// Save CounterPluginState
	store.Set(cp.StateKey(), wire.BinaryBytes(cpState))

	return wrsp.OK
}

func (cp *CounterPlugin) InitChain(store types.KVStore, vals []*wrsp.Validator) {
}

func (cp *CounterPlugin) BeginBlock(store types.KVStore, hash []byte, header *wrsp.Header) {
}

func (cp *CounterPlugin) EndBlock(store types.KVStore, height uint64) (res wrsp.ResponseEndBlock) {
	return
}
