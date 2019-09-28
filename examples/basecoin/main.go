package main

import (
	"fmt"
	"os"

	"github.com/tepleton/wrsp/server"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/app"
	"github.com/tepleton/tepleton-sdk/store"
	"github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/sendtx"
)

func main() {

	app := app.NewApp("basecoin")

	db, err := dbm.NewGoLevelDB("basecoin", "basecoin-data")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create CommitStoreLoader
	cacheSize := 10000
	numHistory := int64(100)
	loader := store.NewIAVLStoreLoader(db, cacheSize, numHistory)

	// Create MultiStore
	multiStore := store.NewCommitMultiStore(db)
	multiStore.SetSubstoreLoader("main", loader)

	// Create Handler
	handler := types.ChainDecorators(
		// recover.Decorator(),
		// logger.Decorator(),
		auth.DecoratorFn(newAccountStore),
	).WithHandler(
		sendtx.TransferHandlerFn(newAccountStore),
	)

	// TODO: load genesis
	// TODO: InitChain with validators
	// accounts := newAccountStore(multiStore.GetKVStore("main"))
	// TODO: set the genesis accounts

	// Set everything on the app and load latest
	app.SetCommitMultiStore(multiStore)
	app.SetTxParser(txParser)
	app.SetHandler(handler)
	if err := app.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	return
}
