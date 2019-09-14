package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/tmlibs/cli"
	tmflags "github.com/tepleton/tmlibs/cli/flags"
	"github.com/tepleton/tmlibs/log"
)

const (
	defaultLogLevel = "error"
	FlagLogLevel    = "log_level"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

var RootCmd = &cobra.Command{
	Use:   "basecoin",
	Short: "A cryptocurrency framework in Golang based on Tendermint-Core",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		level := viper.GetString(FlagLogLevel)
		logger, err = tmflags.ParseLogLevel(level, logger, defaultLogLevel)
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().String(FlagLogLevel, defaultLogLevel, "Log level")
}
