package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/libs/cli"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"
	tmtypes "github.com/tepleton/tepleton/types"

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
	if err != nil {
		return
	}
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

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	dapp := app.NewDemocoinApp(logger, db)
	return dapp.ExportAppStateAndValidators()
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
		server.ConstructAppCreator(newApp, "democoin"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "democoin"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.democoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
