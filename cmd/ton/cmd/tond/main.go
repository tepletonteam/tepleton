package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/cmd/ton/app"
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

func generateApp(rootDir string, logger log.Logger) (wrsp.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("ton", dataDir)
	if err != nil {
		return nil, err
	}
	bapp := app.NewBasecoinApp(logger, db)
	//dbAcc, err := dbm.NewGoLevelDB("ton-acc", dataDir)
	//if err != nil {
	//return nil, err
	//}
	//dbIBC, err := dbm.NewGoLevelDB("ton-ibc", dataDir)
	//if err != nil {
	//return nil, err
	//}
	//dbStake, err := dbm.NewGoLevelDB("ton-stake", dataDir)
	//if err != nil {
	//return nil, err
	//}
	//dbs := map[string]dbm.DB{
	//"main":  dbMain,
	//"acc":   dbAcc,
	//"ibc":   dbIBC,
	//"stake": dbStake,
	//}
	//bapp := app.NewGaiaApp(logger, dbs)
	return bapp, nil
}

func main() {
	server.AddCommands(rootCmd, app.DefaultGenAppState, generateApp, context)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.tond")
	executor := cli.PrepareBaseCmd(rootCmd, "GA", rootDir)
	executor.Execute()
}
