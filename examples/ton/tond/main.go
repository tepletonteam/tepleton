package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wrsp "github.com/tepleton/wrsp/types"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	cfg "github.com/tepleton/tepleton/config"
	"github.com/tepleton/tmlibs/cli"
	tmflags "github.com/tepleton/tmlibs/cli/flags"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/version"
)

// tondCmd is the entry point for this binary
var (
	context  = server.NewContext(nil, nil)
	tondCmd = &cobra.Command{
		Use:   "tond",
		Short: "Gaia Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == version.VersionCmd.Name() {
				return nil
			}
			config, err := tcmd.ParseConfig()
			if err != nil {
				return err
			}
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
			if err != nil {
				return err
			}
			if viper.GetBool(cli.TraceFlag) {
				logger = log.NewTracingLogger(logger)
			}
			logger = logger.With("module", "main")
			context.Config = config
			context.Logger = logger
			return nil
		},
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, string, cmn.HexBytes, error) {
	addr, secret, err := server.GenerateCoinKey()
	if err != nil {
		return nil, "", nil, err
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
	return json.RawMessage(opts), secret, addr, nil
}

func generateApp(rootDir string, logger log.Logger) (wrsp.Application, error) {
	// TODO: set this to something real
	app := new(baseapp.BaseApp)
	return app, nil
}

func main() {
	tondCmd.AddCommand(
		server.InitCmd(defaultOptions, context),
		server.StartCmd(generateApp, context),
		server.UnsafeResetAllCmd(context),
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(tondCmd, "GA", os.ExpandEnv("$HOME/.tond"))
	executor.Execute()
}
