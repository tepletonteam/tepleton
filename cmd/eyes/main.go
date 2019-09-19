package main

import (
	"os"

	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/basecoin"
	client "github.com/tepleton/basecoin/client/commands"
	"github.com/tepleton/basecoin/cmd/basecoin/commands"
	"github.com/tepleton/basecoin/modules/base"
	"github.com/tepleton/basecoin/modules/eyes"
	"github.com/tepleton/basecoin/stack"
)

// BuildApp constructs the stack we want to use for this app
func BuildApp() basecoin.Handler {
	return stack.New(
		base.Logger{},
		stack.Recovery{},
	).
		// We do this to demo real usage, also embeds it under it's own namespace
		Dispatch(
			stack.WrapHandler(eyes.NewHandler()),
		)
}

func main() {
	rt := commands.RootCmd
	rt.Short = "eyes"
	rt.Long = "A demo app to show key-value store with proofs over wrsp"

	commands.Handler = BuildApp()

	rt.AddCommand(
		// out own init command to not require argument
		InitCmd,
		commands.StartCmd,
		commands.UnsafeResetAllCmd,
		client.VersionCmd,
	)

	cmd := cli.PrepareMainCmd(rt, "EYE", os.ExpandEnv("$HOME/.eyes"))
	cmd.Execute()
}
