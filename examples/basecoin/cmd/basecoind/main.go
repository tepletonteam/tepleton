package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tepleton/tmlibs/cli"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/version"
)

// tondCmd is the entry point for this binary
var (
	tondCmd = &cobra.Command{
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

func main() {
	// TODO: this should somehow be updated on cli flags?
	// But we need to create the app first... hmmm.....
	rootDir := os.ExpandEnv("$HOME/.basecoind")

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	db, err := dbm.NewGoLevelDB("basecoin", rootDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bapp := app.NewBasecoinApp(logger, db)

	tondCmd.AddCommand(
		server.InitCmd(defaultOptions, bapp.Logger),
		server.StartCmd(bapp, bapp.Logger),
		server.UnsafeResetAllCmd(bapp.Logger),
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(tondCmd, "BC", rootDir)
	executor.Execute()
}
