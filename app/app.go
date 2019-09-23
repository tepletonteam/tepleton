package app

import (
	"bytes"
	"fmt"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/errors"
	"github.com/tepleton/tepleton-sdk/stack"
	sm "github.com/tepleton/tepleton-sdk/state"
	"github.com/tepleton/tepleton-sdk/version"
)

//nolint
const (
	ModuleNameBase = "base"
	ChainKey       = "chain_id"
)

// Basecoin - The WRSP application
type Basecoin struct {
	*BaseApp
	handler sdk.Handler
	tick    Ticker
}

// Ticker - tick function
type Ticker func(sm.SimpleDB) ([]*wrsp.Validator, error)

var _ wrsp.Application = &Basecoin{}

// NewBasecoin - create a new instance of the basecoin application
func NewBasecoin(handler sdk.Handler, store *Store, logger log.Logger) *Basecoin {
	return &Basecoin{
		BaseApp: NewBaseApp(store, logger),
		handler: handler,
	}
}

// NewBasecoinTick - create a new instance of the basecoin application with tick functionality
func NewBasecoinTick(handler sdk.Handler, store *Store, logger log.Logger, tick Ticker) *Basecoin {
	return &Basecoin{
		BaseApp: NewBaseApp(store, logger),
		handler: handler,
		tick:    tick,
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

	log, err := app.handler.InitState(app.Logger(), state, module, key, value)
	if err == nil {
		return log
	}
	return "Error: " + err.Error()
}

// DeliverTx - WRSP
func (app *Basecoin) DeliverTx(txBytes []byte) wrsp.Result {
	tx, err := sdk.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.height,
		app.Logger().With("call", "delivertx"),
	)
	res, err := app.handler.DeliverTx(ctx, app.state.Append(), tx)

	if err != nil {
		return errors.Result(err)
	}
	app.AddValChange(res.Diff)
	return sdk.ToWRSP(res)
}

// CheckTx - WRSP
func (app *Basecoin) CheckTx(txBytes []byte) wrsp.Result {
	tx, err := sdk.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.height,
		app.Logger().With("call", "checktx"),
	)
	res, err := app.handler.CheckTx(ctx, app.state.Check(), tx)

	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToWRSP(res)
}

// BeginBlock - WRSP
func (app *Basecoin) BeginBlock(req wrsp.RequestBeginBlock) {
	// call the embeded Begin
	app.BaseApp.BeginBlock(req)

	// now execute tick
	if app.tick != nil {
		diff, err := app.tick(app.state.Append())
		if err != nil {
			panic(err)
		}
		app.AddValChange(diff)
	}
}

/////////////////////////// Move to SDK ///////

// BaseApp contains a data store and all info needed
// to perform queries and handshakes.
//
// It should be embeded in another struct for CheckTx,
// DeliverTx and initializing state from the genesis.
type BaseApp struct {
	info  *sm.ChainState
	state *Store

	pending []*wrsp.Validator
	height  uint64
	logger  log.Logger
}

// NewBaseApp creates a data store to handle queries
func NewBaseApp(store *Store, logger log.Logger) *BaseApp {
	return &BaseApp{
		info:   sm.NewChainState(),
		state:  store,
		logger: logger,
	}
}

// GetChainID returns the currently stored chain
func (app *BaseApp) GetChainID() string {
	return app.info.GetChainID(app.state.Committed())
}

// GetState returns the delivertx state, should be removed
func (app *BaseApp) GetState() sm.SimpleDB {
	return app.state.Append()
}

// Logger returns the application base logger
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// Info - WRSP
func (app *BaseApp) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {
	resp := app.state.Info()
	app.logger.Debug("Info",
		"height", resp.LastBlockHeight,
		"hash", fmt.Sprintf("%X", resp.LastBlockAppHash))
	app.height = resp.LastBlockHeight
	return wrsp.ResponseInfo{
		Data:             fmt.Sprintf("Basecoin v%v", version.Version),
		LastBlockHeight:  resp.LastBlockHeight,
		LastBlockAppHash: resp.LastBlockAppHash,
	}
}

// SetOption - WRSP
func (app *BaseApp) SetOption(key string, value string) string {
	return "Not Implemented"
}

// Query - WRSP
func (app *BaseApp) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Log = "Query cannot be zero length"
		resQuery.Code = wrsp.CodeType_EncodingError
		return
	}

	return app.state.Query(reqQuery)
}

// Commit - WRSP
func (app *BaseApp) Commit() (res wrsp.Result) {
	// Commit state
	res = app.state.Commit()
	if res.IsErr() {
		cmn.PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// InitChain - WRSP
func (app *BaseApp) InitChain(req wrsp.RequestInitChain) {
	// for _, plugin := range app.plugins.GetList() {
	// 	plugin.InitChain(app.state, validators)
	// }
}

// BeginBlock - WRSP
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) {
	app.height++
}

// EndBlock - WRSP
// Returns a list of all validator changes made in this block
func (app *BaseApp) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// TODO: cleanup in case a validator exists multiple times in the list
	res.Diffs = app.pending
	app.pending = nil
	return
}

func (app *BaseApp) AddValChange(diffs []*wrsp.Validator) {
	for _, d := range diffs {
		idx := pubKeyIndex(d, app.pending)
		if idx >= 0 {
			app.pending[idx] = d
		} else {
			app.pending = append(app.pending, d)
		}
	}
}

// return index of list with validator of same PubKey, or -1 if no match
func pubKeyIndex(val *wrsp.Validator, list []*wrsp.Validator) int {
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
