package app

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	sdk "github.com/tepleton/tepleton-sdk"
)

// InitApp - The WRSP application with initialization hooks
type InitApp struct {
	*BaseApp
	initState sdk.InitStater
	initVals  sdk.InitValidator
}

var _ wrsp.Application = &InitApp{}

// NewInitApp extends a BaseApp with initialization callbacks,
// which it binds to the proper wrsp calls
func NewInitApp(base *BaseApp, initState sdk.InitStater,
	initVals sdk.InitValidator) *InitApp {

	return &InitApp{
		BaseApp:   base,
		initState: initState,
		initVals:  initVals,
	}
}

// InitState - used to setup state (was SetOption)
// to be call from setting up the genesis file
func (app *InitApp) InitState(module, key, value string) error {
	state := app.Append()
	logger := app.Logger().With("module", module, "key", key)

	if module == sdk.ModuleNameBase {
		if key == sdk.ChainKey {
			app.info.SetChainID(state, value)
			return nil
		}
		logger.Error("Invalid genesis option")
		return fmt.Errorf("Unknown base option: %s", key)
	}

	log, err := app.initState.InitState(logger, state, module, key, value)
	if err != nil {
		logger.Error("Invalid genesis option", "err", err)
	} else {
		logger.Info(log)
	}
	return err
}

// InitChain - WRSP - sets the initial validators
func (app *InitApp) InitChain(req wrsp.RequestInitChain) {
	// return early if no InitValidator registered
	if app.initVals == nil {
		return
	}

	logger, store := app.Logger(), app.Append()
	app.initVals.InitValidators(logger, store, req.Validators)
}
