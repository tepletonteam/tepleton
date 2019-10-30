package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	"github.com/tepleton/tepleton-sdk/server"
)

// rootCmd is the entry point for this binary
var (
	context = server.NewDefaultContext()
	rootCmd = &cobra.Command{
		Use:               "tond",
		Short:             "Gaia Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(context),
	}
)

// TODO: distinguish from basecoin
func generateApp(rootDir string, logger log.Logger) (wrsp.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("ton", dataDir)
	if err != nil {
		return nil, err
	}
	bapp := app.NewBasecoinApp(logger, db)
	return bapp, nil
}

func main() {
	server.AddCommands(rootCmd, server.DefaultGenAppState, generateApp, context)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", os.ExpandEnv("$HOME/.tond"))
	executor.Execute()
}
