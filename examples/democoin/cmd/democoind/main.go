package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/examples/democoin/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/wire"
)

// init parameters
var CoolAppInit = server.AppInit{
	AppGenState: CoolAppGenState,
	AppGenTx:    server.SimpleAppGenTx,
}

// coolGenAppParams sets up the app_state and appends the cool app state
func CoolAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	appState, err = server.SimpleAppGenState(cdc, appGenTxs)
	if err != nil {
		return
	}
	key := "cool"
	value := json.RawMessage(`{
        "trend": "ice-cold"
      }`)
	appState, err = server.AppendJSON(cdc, appState, key, value)
	key = "pow"
	value = json.RawMessage(`{
        "difficulty": 1,
        "count": 0
      }`)
	appState, err = server.AppendJSON(cdc, appState, key, value)
	return
}

func newApp(logger log.Logger, db dbm.DB) wrsp.Application {
	return app.NewDemocoinApp(logger, db)
}

func exportApp(logger log.Logger, db dbm.DB) (interface{}, *wire.Codec) {
	dapp := app.NewDemocoinApp(logger, db)
	return dapp.ExportGenesis(), app.MakeCodec()
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "democoind",
		Short:             "Democoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, CoolAppInit,
		baseapp.GenerateFn(newApp, "democoin"),
		baseapp.ExportFn(exportApp, "democoin"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.democoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	executor.Execute()
}
