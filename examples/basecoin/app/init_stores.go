package app

import (
	"github.com/tepleton/tepleton-sdk/examples/basecoin/types"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *BasecoinApp) initStores() {
	app.mountStores()
	app.initAccountMapper()
}

// Initialize root stores.
func (app *BasecoinApp) mountStores() {

	// Create MultiStore mounts.
	app.BaseApp.MountStore(app.capKeyMainStore, sdk.StoreTypeIAVL)
	app.BaseApp.MountStore(app.capKeyIBCStore, sdk.StoreTypeIAVL)
}

// Initialize the AccountMapper.
func (app *BasecoinApp) initAccountMapper() {

	var accountMapper = auth.NewAccountMapper(
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)

	// Register all interfaces and concrete types that
	// implement those interfaces, here.
	cdc := accountMapper.WireCodec()
	auth.RegisterWireBaseAccount(cdc)

	// Make accountMapper's WireCodec() inaccessible.
	app.accountMapper = accountMapper.Seal()
}
