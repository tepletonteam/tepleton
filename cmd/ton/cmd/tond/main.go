package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	tmtypes "github.com/tepleton/tepleton/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

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
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) wrsp.Application {
	return app.NewGaiaApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gapp := app.NewGaiaApp(logger, db)
	return gapp.ExportAppStateAndValidators()
}
