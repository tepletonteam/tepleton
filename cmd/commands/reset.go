package commands

import (
	"github.com/spf13/cobra"

	tmcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	tmcfg "github.com/tepleton/tepleton/config/tepleton"
)

var UnsafeResetAllCmd = &cobra.Command{
	Use:   "unsafe_reset_all",
	Short: "Reset all blockchain data",
	RunE:  unsafeResetAllCmd,
}

func unsafeResetAllCmd(cmd *cobra.Command, args []string) error {
	basecoinDir := BasecoinRoot("")
	tmConfig := tmcfg.GetConfig(basecoinDir)
	tmcmd.ResetAll(tmConfig, log)
	return nil
}
