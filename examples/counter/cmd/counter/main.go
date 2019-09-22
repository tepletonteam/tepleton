package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tepleton/tmlibs/cli"

	client "github.com/tepleton/tepleton-sdk/client/commands"
	"github.com/tepleton/tepleton-sdk/examples/counter/plugins/counter"
	"github.com/tepleton/tepleton-sdk/server/commands"
)

// RootCmd is the entry point for this binary
var RootCmd = &cobra.Command{
	Use:   "counter",
	Short: "demo application for tepleton sdk",
}

func main() {

	// TODO: register the counter here
	commands.Handler = counter.NewHandler("strings")

	RootCmd.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		commands.UnsafeResetAllCmd,
		client.VersionCmd,
	)
	commands.SetUpRoot(RootCmd)

	cmd := cli.PrepareMainCmd(RootCmd, "CT", os.ExpandEnv("$HOME/.counter"))
	cmd.Execute()
}
