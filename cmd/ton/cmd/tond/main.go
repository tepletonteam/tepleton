package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/libs/cli"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"
	tmtypes "github.com/tepleton/tepleton/types"

	"github.com/tepleton/tepleton-sdk/cmd/ton/app"
	"github.com/tepleton/tepleton-sdk/server"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "tond",
		Short:             "Gaia Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.GaiaAppInit(),
		server.ConstructAppCreator(newApp, "ton"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "ton"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB) wrsp.Application {
	return app.NewGaiaApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gapp := app.NewGaiaApp(logger, db)
	return gapp.ExportAppStateAndValidators()
}
