package app

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/errors"
	"github.com/tepleton/tepleton-sdk/stack"
)

// BaseApp - The WRSP application
type BaseApp struct {
	*StoreApp
	handler sdk.Handler
	clock   sdk.Ticker
}

var _ wrsp.Application = &BaseApp{}

// NewBaseApp extends a StoreApp with a handler and a ticker,
// which it binds to the proper wrsp calls
func NewBaseApp(store *StoreApp, handler sdk.Handler, clock sdk.Ticker) *BaseApp {
	return &BaseApp{
		StoreApp: store,
		handler:  handler,
		clock:    clock,
	}
}

// DeliverTx - WRSP - dispatches to the handler
func (app *BaseApp) DeliverTx(txBytes []byte) wrsp.Result {
	tx, err := sdk.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.WorkingHeight(),
		app.Logger().With("call", "delivertx"),
	)
	res, err := app.handler.DeliverTx(ctx, app.Append(), tx)

	if err != nil {
		return errors.Result(err)
	}
	app.AddValChange(res.Diff)
	return sdk.ToWRSP(res)
}

// CheckTx - WRSP - dispatches to the handler
func (app *BaseApp) CheckTx(txBytes []byte) wrsp.Result {
	tx, err := sdk.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.WorkingHeight(),
		app.Logger().With("call", "checktx"),
	)
	res, err := app.handler.CheckTx(ctx, app.Check(), tx)

	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToWRSP(res)
}

// BeginBlock - WRSP - triggers Tick actions
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) {
	// execute tick if present
	if app.clock != nil {
		ctx := stack.NewContext(
			app.GetChainID(),
			app.WorkingHeight(),
			app.Logger().With("call", "tick"),
		)

		diff, err := app.clock.Tick(ctx, app.Append())
		if err != nil {
			panic(err)
		}
		app.AddValChange(diff)
	}
}

// InitState - used to setup state (was SetOption)
// to be used by InitChain later
//
// TODO: rethink this a bit more....
func (app *BaseApp) InitState(module, key, value string) (string, error) {
	state := app.Append()

	if module == sdk.ModuleNameBase {
		if key == sdk.ChainKey {
			app.info.SetChainID(state, value)
			return "Success", nil
		}
		return "", fmt.Errorf("unknown base option: %s", key)
	}

	log, err := app.handler.InitState(app.Logger(), state, module, key, value)
	if err != nil {
		app.Logger().Error("Genesis App Options", "err", err)
	} else {
		app.Logger().Info(log)
	}
	return log, err
}
