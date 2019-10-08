package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/version"
)

// basecoindCmd is the entry point for this binary
var (
	basecoindCmd = &cobra.Command{
		Use:   "tond",
		Short: "Gaia Daemon (server)",
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, error) {
	addr, secret, err := server.GenerateCoinKey()
	if err != nil {
		return nil, err
	}
	fmt.Println("Secret phrase to access coins:")
	fmt.Println(secret)

	opts := fmt.Sprintf(`{
      "accounts": [{
        "address": "%s",
        "coins": [
          {
            "denom": "mycoin",
            "amount": 9007199254740992
          }
        ]
      }]
    }`, addr)
	return json.RawMessage(opts), nil
}

func generateApp(rootDir string, logger log.Logger) wrsp.Application {
	db, err := dbm.NewGoLevelDB("basecoin", rootDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bapp := app.NewBasecoinApp(logger, db)
	return bapp
}

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")

	basecoindCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		version.VersionCmd,
	)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.basecoind")
	executor := cli.PrepareBaseCmd(basecoindCmd, "BC", rootDir)
	executor.Execute()
}
