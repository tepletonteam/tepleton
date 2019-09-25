package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tepleton/tmlibs/cli"

	sdk "github.com/tepleton/tepleton-sdk"
	client "github.com/tepleton/tepleton-sdk/client/commands"
	eyesmod "github.com/tepleton/tepleton-sdk/modules/eyes"
	"github.com/tepleton/tepleton-sdk/server/commands"
	"github.com/tepleton/tepleton-sdk/util"

	"github.com/tepleton/tepleton-sdk/examples/eyes"
)

// RootCmd is the entry point for this binary
var RootCmd = &cobra.Command{
	Use:   "eyes",
	Short: "key-value store",
	Long:  "A demo app to show key-value store with proofs over wrsp",
}

// BuildApp constructs the stack we want to use for this app
func BuildApp() sdk.Handler {
	return sdk.ChainDecorators(
		util.Logger{},
		util.Recovery{},
		eyes.Parser{},
		util.Chain{},
	).WithHandler(
		eyesmod.NewHandler(),
	)
}

func main() {
	commands.Handler = BuildApp()

	RootCmd.AddCommand(
		// out own init command to not require argument
		InitCmd,
		commands.StartCmd,
		commands.UnsafeResetAllCmd,
		client.VersionCmd,
	)
	commands.SetUpRoot(RootCmd)

	cmd := cli.PrepareMainCmd(RootCmd, "EYE", os.ExpandEnv("$HOME/.eyes"))
	cmd.Execute()
}
