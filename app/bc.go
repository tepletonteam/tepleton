package app

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/log"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/errors"
	"github.com/tepleton/tepleton-sdk/stack"
	sm "github.com/tepleton/tepleton-sdk/state"
	"github.com/tepleton/tepleton-sdk/version"
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
func NewBasecoin(handler sdk.Handler, dbName string, cacheSize int, logger log.Logger) (*Basecoin, error) {
	appName := fmt.Sprintf("Basecoin v%v", version.Version)
	base, err := NewBaseApp(appName, dbName, cacheSize, logger)
	app := &Basecoin{
		BaseApp: base,
		handler: handler,
	}
	return app, err
}

// NewBasecoinTick - create a new instance of the basecoin application with tick functionality
func NewBasecoinTick(handler sdk.Handler, tick Ticker, dbName string, cacheSize int, logger log.Logger) (*Basecoin, error) {
	appName := fmt.Sprintf("Basecoin v%v", version.Version)
	base, err := NewBaseApp(appName, dbName, cacheSize, logger)
	app := &Basecoin{
		BaseApp: base,
		handler: handler,
		tick:    tick,
	}
	return app, err
}

// InitState - used to setup state (was SetOption)
// to be used by InitChain later
func (app *Basecoin) InitState(module, key, value string) (string, error) {
	state := app.Append()

	if module == ModuleNameBase {
		if key == ChainKey {
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

// DeliverTx - WRSP
func (app *Basecoin) DeliverTx(txBytes []byte) wrsp.Result {
	tx, err := sdk.LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	ctx := stack.NewContext(
		app.GetChainID(),
		app.height+1,
		app.Logger().With("call", "delivertx"),
	)
	res, err := app.handler.DeliverTx(ctx, app.Append(), tx)

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
		app.height+1,
		app.Logger().With("call", "checktx"),
	)
	res, err := app.handler.CheckTx(ctx, app.Check(), tx)

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
		diff, err := app.tick(app.Append())
		if err != nil {
			panic(err)
		}
		app.AddValChange(diff)
	}
}
