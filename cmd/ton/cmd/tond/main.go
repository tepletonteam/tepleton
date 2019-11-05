package main

import (
	"path/filepath"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

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

	server.AddCommands(ctx, cdc, rootCmd, app.GaiaAppInit(), generateApp)
	server.AddCommands(ctx, cdc, rootCmd, app.GaiaAppInit(), generateApp, exportApp)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	executor.Execute()
}

func generateApp(rootDir string, logger log.Logger) (wrsp.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("ton", dataDir)
	if err != nil {
		return nil, err
	}
	bapp := app.NewGaiaApp(logger, db)
	return bapp, nil
}

func exportApp(rootDir string, logger log.Logger) (interface{}, *wire.Codec, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("ton", dataDir)
	if err != nil {
		return nil, nil, err
	}
	bapp := app.NewGaiaApp(log.NewNopLogger(), db)
	if err != nil {
		return nil, nil, err
	}
	return bapp.ExportGenesis(), app.MakeCodec(), nil
}
