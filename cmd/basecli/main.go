package main

import (
	"os"

	"github.com/spf13/cobra"
	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/light-client/commands"
	"github.com/tepleton/light-client/commands/proofs"
	"github.com/tepleton/light-client/commands/seeds"
	"github.com/tepleton/light-client/commands/txs"
	"github.com/tepleton/tmlibs/cli"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "basecli",
	Short: "Light client for tepleton",
	Long: `Basecli is an version of tmcli including custom logic to
present a nice (not raw hex) interface to the basecoin blockchain structure.

This is a useful tool, but also serves to demonstrate how one can configure
tmcli to work for any custom wrsp app.
`,
}

func main() {
	commands.AddBasicFlags(BaseCli)

	//initialize proofs and txs
	proofs.StatePresenters.Register("account", AccountPresenter{})
	proofs.TxPresenters.Register("base", BaseTxPresenter{})
	txs.Register("send", SendTxMaker{})

	// set up the various commands to use
	BaseCli.AddCommand(
		keycmd.RootCmd,
		commands.InitCmd,
		seeds.RootCmd,
		proofs.RootCmd,
		txs.RootCmd,
	)

	cmd := cli.PrepareMainCmd(BaseCli, "BC", os.ExpandEnv("$HOME/.basecli"))
	cmd.Execute()
}
