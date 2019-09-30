package app

import (
	apm "github.com/tepleton/tepleton-sdk/app"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
)

// initSDKApp() happens after initCapKeys() and initStores().
// initSDKApp() happens before initRoutes().
func (app *BasecoinApp) initSDKApp() {
	app.App = apm.NewApp(appName, app.multiStore)
	app.initSDKAppTxDecoder()
	app.initSDKAppAnteHandler()
}

func (app *BasecoinApp) initSDKAppTxDecoder() {
	cdc := makeTxCodec()
	app.App.SetTxDecoder(func(txBytes []byte) (sdk.Tx, error) {
		var tx = sdk.StdTx{}
		err := cdc.UnmarshalBinary(txBytes, &tx)
		return tx, err
	})
}

func (app *BasecoinApp) initSDKAppAnteHandler() {
	var authAnteHandler = auth.NewAnteHandler(app.accStore)
	app.App.SetDefaultAnteHandler(authAnteHandler)
}
