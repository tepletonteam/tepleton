package server

import (
	"github.com/spf13/cobra"

	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tmlibs/log"
)

// UnsafeResetAllCmd - extension of the tepleton command, resets initialization
func UnsafeResetAllCmd(logger log.Logger) *cobra.Command {
	cmd := resetAll{logger}
	return &cobra.Command{
		Use:   "unsafe_reset_all",
		Short: "Reset all blockchain data",
		RunE:  cmd.run,
	}
}

type resetAll struct {
	logger log.Logger
}

func (r resetAll) run(cmd *cobra.Command, args []string) error {
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}
	tcmd.ResetAll(cfg.DBDir(), cfg.PrivValidatorFile(), r.logger)
	return nil
}
