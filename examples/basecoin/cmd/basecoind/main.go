package main

import (
	"encoding/json"
	"os"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/spf13/cobra"
	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/libs/cli"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"
	tmtypes "github.com/tepleton/tepleton/types"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "basecoind",
		Short:             "Basecoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "basecoin"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "basecoin"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.basecoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB) wrsp.Application {
	return app.NewBasecoinApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewBasecoinApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}
