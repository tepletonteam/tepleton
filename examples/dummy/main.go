package main

import (
	"fmt"
	"os"

	"github.com/tepleton/wrsp/server"
	cmn "github.com/tepleton/tmlibs/common"

	bam "github.com/tepleton/tepleton-sdk/baseapp"
	sdk "github.com/tepleton/tepleton-sdk/types"
)

func main() {

	// Capabilities key to access the main KVStore.
	var capKeyMainStore = sdk.NewKVStoreKey("main")

	// Create BaseApp.
	var baseApp = bam.NewBaseApp("dummy")

	// Set mounts for BaseApp's MultiStore.
	baseApp.MountStore(capKeyMainStore, sdk.StoreTypeIAVL)

	// Set Tx decoder
	baseApp.SetTxDecoder(decodeTx)

	// Set a handler Route.
	baseApp.Router().AddRoute("dummy", DummyHandler(capKeyMainStore))

	// Load latest version.
	if err := baseApp.LoadLatestVersion(capKeyMainStore); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start the WRSP server
	srv, err := server.NewServer("0.0.0.0:46658", "socket", baseApp)
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

func DummyHandler(storeKey sdk.StoreKey) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		// tx is already unmarshalled
		key := tx.Get("key").([]byte)
		value := tx.Get("value").([]byte)

		store := ctx.KVStore(storeKey)
		store.Set(key, value)

		return sdk.Result{
			Code: 0,
			Log:  fmt.Sprintf("set %s=%s", key, value),
		}
	}
}
