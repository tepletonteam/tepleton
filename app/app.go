package app

import (
	"fmt"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin"
	eyes "github.com/tepleton/merkleeyes/client"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/stack"
	sm "github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/version"
)

const (
	ModuleNameBase = "base"
	ChainKey       = "chain_id"
)

type Basecoin struct {
	eyesCli    *eyes.Client
	state      *sm.State
	cacheState *sm.State
	handler    basecoin.Handler
	logger     log.Logger
}

func NewBasecoin(h basecoin.Handler, eyesCli *eyes.Client, l log.Logger) *Basecoin {
	state := sm.NewState(eyesCli, l.With("module", "state"))

	return &Basecoin{
		handler:    h,
		eyesCli:    eyesCli,
		state:      state,
		cacheState: nil,
		logger:     l,
	}
}

// placeholder to just handle sendtx
func DefaultHandler() basecoin.Handler {
	// use the default stack
	h := coin.NewHandler()
	d := stack.NewDispatcher(stack.WrapHandler(h))
	return stack.NewDefault().Use(d)
}

// XXX For testing, not thread safe!
func (app *Basecoin) GetState() *sm.State {
	return app.state.CacheWrap()
}

// WRSP::Info
func (app *Basecoin) Info() wrsp.ResponseInfo {
	resp, err := app.eyesCli.InfoSync()
	if err != nil {
		cmn.PanicCrisis(err)
	}
	return wrsp.ResponseInfo{
		Data:             cmn.Fmt("Basecoin v%v", version.Version),
		LastBlockHeight:  resp.LastBlockHeight,
		LastBlockAppHash: resp.LastBlockAppHash,
	}
}

// WRSP::SetOption
func (app *Basecoin) SetOption(key string, value string) string {
	module, prefix := splitKey(key)
	if module == ModuleNameBase {
		return app.setBaseOption(prefix, value)
	}

	log, err := app.handler.SetOption(app.logger, app.state, module, prefix, value)
	if err == nil {
		return log
	}
	return "Error: " + err.Error()
}

func (app *Basecoin) setBaseOption(key, value string) string {
	if key == ChainKey {
		app.state.SetChainID(value)
		return "Success"
	}
	return fmt.Sprintf("Error: unknown base option: %s", key)
}

// WRSP::DeliverTx
func (app *Basecoin) DeliverTx(txBytes []byte) wrsp.Result {
	tx, err := basecoin.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	// TODO: can we abstract this setup and commit logic??
	cache := app.state.CacheWrap()
	ctx := stack.NewContext(
		app.state.GetChainID(),
		app.logger.With("call", "delivertx"),
	)
	res, err := app.handler.DeliverTx(ctx, cache, tx)

	if err != nil {
		// discard the cache...
		return errors.Result(err)
	}
	// commit the cache and return result
	cache.CacheSync()
	return res.ToWRSP()
}

// WRSP::CheckTx
func (app *Basecoin) CheckTx(txBytes []byte) wrsp.Result {
	tx, err := basecoin.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	// TODO: can we abstract this setup and commit logic??
	ctx := stack.NewContext(
		app.state.GetChainID(),
		app.logger.With("call", "checktx"),
	)
	// checktx generally shouldn't touch the state, but we don't care
	// here on the framework level, since the cacheState is thrown away next block
	res, err := app.handler.CheckTx(ctx, app.cacheState, tx)

	if err != nil {
		return errors.Result(err)
	}
	return res.ToWRSP()
}

// WRSP::Query
func (app *Basecoin) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Log = "Query cannot be zero length"
		resQuery.Code = wrsp.CodeType_EncodingError
		return
	}

	resQuery, err := app.eyesCli.QuerySync(reqQuery)
	if err != nil {
		resQuery.Log = "Failed to query MerkleEyes: " + err.Error()
		resQuery.Code = wrsp.CodeType_InternalError
		return
	}
	return
}

// WRSP::Commit
func (app *Basecoin) Commit() (res wrsp.Result) {

	// Commit state
	res = app.state.Commit()

	// Wrap the committed state in cache for CheckTx
	app.cacheState = app.state.CacheWrap()

	if res.IsErr() {
		cmn.PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// WRSP::InitChain
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	// for _, plugin := range app.plugins.GetList() {
	// 	plugin.InitChain(app.state, validators)
	// }
}

// WRSP::BeginBlock
func (app *Basecoin) BeginBlock(hash []byte, header *wrsp.Header) {
	// for _, plugin := range app.plugins.GetList() {
	// 	plugin.BeginBlock(app.state, hash, header)
	// }
}

// WRSP::EndBlock
func (app *Basecoin) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// for _, plugin := range app.plugins.GetList() {
	// 	pluginRes := plugin.EndBlock(app.state, height)
	// 	res.Diffs = append(res.Diffs, pluginRes.Diffs...)
	// }
	return
}

// Splits the string at the first '/'.
// if there are none, assign default module ("base").
func splitKey(key string) (string, string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return ModuleNameBase, key
}
