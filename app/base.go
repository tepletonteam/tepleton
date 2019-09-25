package app

import (
	wrsp "github.com/tepleton/wrsp/types"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/errors"
	"github.com/tepleton/tepleton-sdk/util"
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

	// TODO: use real context on refactor
	ctx := util.MockContext(
		app.GetChainID(),
		app.WorkingHeight(),
	)
	// Note: first decorator must parse bytes
	res, err := app.handler.DeliverTx(ctx, app.Append(), txBytes)

	if err != nil {
		return errors.Result(err)
	}
	app.AddValChange(res.Diff)
	return sdk.ToWRSP(res)
}

// CheckTx - WRSP - dispatches to the handler
func (app *BaseApp) CheckTx(txBytes []byte) wrsp.Result {
	// TODO: use real context on refactor
	ctx := util.MockContext(
		app.GetChainID(),
		app.WorkingHeight(),
	)
	// Note: first decorator must parse bytes
	res, err := app.handler.CheckTx(ctx, app.Check(), txBytes)

	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToWRSP(res)
}

// BeginBlock - WRSP - triggers Tick actions
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) {
	// execute tick if present
	if app.clock != nil {
		ctx := util.MockContext(
			app.GetChainID(),
			app.WorkingHeight(),
		)

		diff, err := app.clock.Tick(ctx, app.Append())
		if err != nil {
			panic(err)
		}
		app.AddValChange(diff)
	}
}
