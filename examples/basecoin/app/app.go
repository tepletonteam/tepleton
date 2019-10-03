package app

import (
	"fmt"
	"os"

	bam "github.com/tepleton/tepleton-sdk/baseapp"
	sdk "github.com/tepleton/tepleton-sdk/types"

	"github.com/tepleton/wrsp/server"
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-wire"
	cmn "github.com/tepleton/tmlibs/common"
)

const appName = "BasecoinApp"

// BasecoinApp - extended WRSP application
type BasecoinApp struct {
	*bam.BaseApp
	router     bam.Router
	cdc        *wire.Codec
	multiStore sdk.CommitMultiStore //TODO distinguish this store from *bam.BaseApp.cms <- is this one master?? confused

	// The key to access the substores.
	capKeyMainStore *sdk.KVStoreKey
	capKeyIBCStore  *sdk.KVStoreKey

	// Object mappers:
	accountMapper sdk.AccountMapper
}

// NewBasecoinApp - create new BasecoinApp
// TODO: This should take in more configuration options.
// TODO: This should be moved into baseapp to isolate complexity
func NewBasecoinApp(genesisPath string) *BasecoinApp {

	// Create and configure app.
	var app = &BasecoinApp{}

	// TODO open up out of functions, or introduce clarity,
	// interdependancies are a nightmare to debug
	app.initCapKeys() // ./init_capkeys.go
	app.initBaseApp() // ./init_baseapp.go
	app.initStores()  // ./init_stores.go
	app.initBaseAppInitStater()
	app.initHandlers() // ./init_handlers.go

	genesisiDoc, err := bam.GenesisDocFromFile(genesisPath)
	if err != nil {
		panic(fmt.Errorf("error loading genesis state: %v", err))
	}

	// TODO: InitChain with validators from genesis transaction bytes

	// very first begin block used for context when setting genesis accounts
	header := wrsp.Header{
		ChainID:        "",
		Height:         0,
		Time:           -1,
		NumTxs:         -1,
		LastCommitHash: []byte{0x00},
		DataHash:       nil,
		ValidatorsHash: nil,
		AppHash:        nil,
	}
	app.BaseApp.BeginBlock(wrsp.RequestBeginBlock{
		Hash:                nil,
		Header:              header,
		AbsentValidators:    nil,
		ByzantineValidators: nil,
	})

	ctxCheckTx := app.BaseApp.NewContext(true, nil)
	ctxDeliverTx := app.BaseApp.NewContext(false, nil)

	err = app.BaseApp.InitStater(ctxCheckTx, ctxDeliverTx, genesisiDoc.AppState)
	if err != nil {
		panic(fmt.Errorf("error loading application genesis state: %v", err))
	}

	app.loadStores()

	return app
}

// RunForever - BasecoinApp execution and cleanup
func (app *BasecoinApp) RunForever() {

	// Start the WRSP server
	srv, err := server.NewServer("0.0.0.0:46658", "socket", app)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srv.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}

// Load the stores
func (app *BasecoinApp) loadStores() {
	if err := app.LoadLatestVersion(app.capKeyMainStore); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
