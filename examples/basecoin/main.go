package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tepleton/wrsp/server"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/app"
	"github.com/tepleton/tepleton-sdk/store"
	"github.com/tepleton/tepleton-sdk/types"
	acm "github.com/tepleton/tepleton-sdk/x/account"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/coinstore"
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
		auth.DecoratorFn(acm.NewAccountStore),
	).WithHandler(
		coinstore.TransferHandlerFn(acm.NewAccountStore),
	)

	// TODO: load genesis
	// TODO: InitChain with validators
	// accounts := acm.NewAccountStore(multiStore.GetKVStore("main"))
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

func txParser(txBytes []byte) (types.Tx, error) {
	var tx coinstore.SendTx
	err := json.Unmarshal(txBytes, &tx)
	return tx, err
}

//-----------------------------------------------------------------------------
