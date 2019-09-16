package app

import (
	"fmt"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/modules/auth"
	"github.com/tepleton/basecoin/modules/base"
	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/modules/fee"
	"github.com/tepleton/basecoin/modules/nonce"
	"github.com/tepleton/basecoin/modules/roles"
	"github.com/tepleton/basecoin/stack"
	sm "github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/state/merkle"
	"github.com/tepleton/basecoin/version"
)

//nolint
const (
	ModuleNameBase = "base"
	ChainKey       = "chain_id"
)

// Basecoin - The WRSP application
type Basecoin struct {
	info *sm.ChainState

	state *merkle.Store

	handler basecoin.Handler
	height  uint64
	logger  log.Logger
}

var _ wrsp.Application = &Basecoin{}

// NewBasecoin - create a new instance of the basecoin application
func NewBasecoin(handler basecoin.Handler, store *merkle.Store, logger log.Logger) *Basecoin {
	return &Basecoin{
		handler: handler,
		info:    sm.NewChainState(),
		state:   store,
		logger:  logger,
	}
}

// DefaultHandler - placeholder to just handle sendtx
func DefaultHandler(feeDenom string) basecoin.Handler {
	// use the default stack
	c := coin.NewHandler()
	r := roles.NewHandler()
	d := stack.NewDispatcher(
		stack.WrapHandler(c),
		stack.WrapHandler(r),
	)
	return stack.New(
		base.Logger{},
		stack.Recovery{},
		auth.Signatures{},
		base.Chain{},
		nonce.ReplayCheck{},
		roles.NewMiddleware(),
		fee.NewSimpleFeeMiddleware(coin.Coin{feeDenom, 0}, fee.Bank),
		stack.Checkpoint{},
	).Use(d)
}

// GetChainID returns the currently stored chain
func (app *Basecoin) GetChainID() string {
	return app.info.GetChainID(app.state.Committed())
}

// GetState is back... please kill me
func (app *Basecoin) GetState() sm.KVStore {
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

// SetOption - WRSP
func (app *Basecoin) SetOption(key string, value string) string {

	module, key := splitKey(key)
	state := app.state.Append()

	if module == ModuleNameBase {
		if key == ChainKey {
			app.info.SetChainID(state, value)
			return "Success"
		}
		return fmt.Sprintf("Error: unknown base option: %s", key)
	}

	log, err := app.handler.SetOption(app.logger, state, module, key, value)
	if err == nil {
		return log
	}
	return "Error: " + err.Error()
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
		// discard the cache...
		return errors.Result(err)
	}
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
func (app *Basecoin) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// for _, plugin := range app.plugins.GetList() {
	// 	pluginRes := plugin.EndBlock(app.state, height)
	// 	res.Diffs = append(res.Diffs, pluginRes.Diffs...)
	// }
	return
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
