package main

import (
	"os"

	"github.com/tepleton/basecoin/app"
	"github.com/tepleton/basecoin/cmd/basecoin/commands"
	"github.com/tepleton/tmlibs/cli"
)

func main() {
	rt := commands.RootCmd

	// require all fees in mycoin - change this in your app!
	commands.Handler = app.DefaultHandler("mycoin")

	rt.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		//commands.RelayCmd,
		commands.UnsafeResetAllCmd,
		commands.VersionCmd,
	)

	cmd := cli.PrepareMainCmd(rt, "BC", os.ExpandEnv("$HOME/.basecoin"))
	cmd.Execute()
}
