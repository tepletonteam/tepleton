package main

import (
	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/cmd/ton/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/wire"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use:               "tond",
		Short:             "Gaia Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.GaiaAppInit(),
		baseapp.GenerateFn(newApp, "ton"),
		baseapp.ExportFn(exportApp, "ton"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) wrsp.Application {
	return app.NewGaiaApp(logger, db)
}

func exportApp(logger log.Logger, db dbm.DB) (interface{}, *wire.Codec) {
	gapp := app.NewGaiaApp(logger, db)
	return gapp.ExportGenesis(), app.MakeCodec()
}
