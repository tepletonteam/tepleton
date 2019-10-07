package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/version"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagFee    = "fee"
)

// tondCmd is the entry point for this binary
var (
	tondCmd = &cobra.Command{
		Use:   "tond",
		Short: "Gaia Daemon (server)",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// TODO: set this to something real
	var node baseapp.BaseApp

	AddNodeCommands(tondCmd, node)
	tondCmd.AddCommand(
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(tondCmd, "GA", os.ExpandEnv("$HOME/.tond"))
	executor.Execute()
}
