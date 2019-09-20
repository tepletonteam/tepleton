package app

import (
	"bytes"
	"fmt"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/stack"
	sm "github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/version"
)

//nolint
const (
	ModuleNameBase = "base"
	ChainKey       = "chain_id"
)

// Basecoin - The WRSP application
type Basecoin struct {
	info  *sm.ChainState
	state *Store

	handler basecoin.Handler

	pending []*wrsp.Validator
	height  uint64
	logger  log.Logger
}

var _ wrsp.Application = &Basecoin{}

// NewBasecoin - create a new instance of the basecoin application
func NewBasecoin(handler basecoin.Handler, store *Store, logger log.Logger) *Basecoin {
	return &Basecoin{
		handler: handler,
		info:    sm.NewChainState(),
		state:   store,
		logger:  logger,
	}
}

// GetChainID returns the currently stored chain
func (app *Basecoin) GetChainID() string {
	return app.info.GetChainID(app.state.Committed())
}

// GetState is back... please kill me
func (app *Basecoin) GetState() sm.SimpleDB {
	return app.state.Append()
}

// Info - WRSP
func (app *Basecoin) Info() wrsp.ResponseInfo {
	resp := app.state.Info()
	app.height = resp.LastBlockHeight
	return wrsp.ResponseInfo{
		Data:             fmt.Sprintf("Basecoin v%v", version.Version),
		LastBlockHeight:  resp.LastBlockHeight,
		LastBlockAppHash: resp.LastBlockAppHash,
	}
}

// InitState - used to setup state (was SetOption)
// to be used by InitChain later
func (app *Basecoin) InitState(key string, value string) string {

	module, key := splitKey(key)
	state := app.state.Append()

	if module == ModuleNameBase {
		if key == ChainKey {
			app.info.SetChainID(state, value)
			return "Success"
		}
		return fmt.Sprintf("Error: unknown base option: %s", key)
	}

	log, err := app.handler.InitState(app.logger, state, module, key, value)
	if err == nil {
		return log
	}
	return "Error: " + err.Error()
}

// SetOption - WRSP
func (app *Basecoin) SetOption(key string, value string) string {
	return "Not Implemented"
}

// DeliverTx - WRSP
func (app *Basecoin) DeliverTx(txBytes []byte) wrsp.Result {
	tx, err := basecoin.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.height,
		app.logger.With("call", "delivertx"),
	)
	res, err := app.handler.DeliverTx(ctx, app.state.Append(), tx)

	if err != nil {
		return errors.Result(err)
	}
	app.addValChange(res.Diff)
	return res.ToWRSP()
}

// CheckTx - WRSP
func (app *Basecoin) CheckTx(txBytes []byte) wrsp.Result {
	tx, err := basecoin.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.height,
		app.logger.With("call", "checktx"),
	)
	res, err := app.handler.CheckTx(ctx, app.state.Check(), tx)

	if err != nil {
		return errors.Result(err)
	}
	return res.ToWRSP()
}

// Query - WRSP
func (app *Basecoin) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Log = "Query cannot be zero length"
		resQuery.Code = wrsp.CodeType_EncodingError
		return
	}

	return app.state.Query(reqQuery)
}

// Commit - WRSP
func (app *Basecoin) Commit() (res wrsp.Result) {
	// Commit state
	res = app.state.Commit()
	if res.IsErr() {
		cmn.PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// InitChain - WRSP
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	// for _, plugin := range app.plugins.GetList() {
	// 	plugin.InitChain(app.state, validators)
	// }
}

// BeginBlock - WRSP
func (app *Basecoin) BeginBlock(hash []byte, header *wrsp.Header) {
	app.height++
	// for _, plugin := range app.plugins.GetList() {
	// 	plugin.BeginBlock(app.state, hash, header)
	// }
}

// EndBlock - WRSP
// Returns a list of all validator changes made in this block
func (app *Basecoin) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// TODO: cleanup in case a validator exists multiple times in the list
	res.Diffs = app.pending
	app.pending = nil
	return
}

func (app *Basecoin) addValChange(diffs []*wrsp.Validator) {
	for _, d := range diffs {
		idx := findVal(d, app.pending)
		if idx >= 0 {
			app.pending[idx] = d
		} else {
			app.pending = append(app.pending, d)
		}
	}
}

// return index of list with validator of same PubKey, or -1 if no match
func findVal(val *wrsp.Validator, list []*wrsp.Validator) int {
	for i, v := range list {
		if bytes.Equal(val.PubKey, v.PubKey) {
			return i
		}
	}
	return -1
}

//TODO move split key to tmlibs?

// Splits the string at the first '/'.
// if there are none, assign default module ("base").
func splitKey(key string) (string, string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return ModuleNameBase, key
}
