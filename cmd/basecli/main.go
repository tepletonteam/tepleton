package main

import (
	"os"

	"github.com/spf13/cobra"
	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/light-client/commands"
	"github.com/tepleton/light-client/commands/proofs"
	"github.com/tepleton/light-client/commands/proxy"
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

func init() {
	commands.AddBasicFlags(BaseCli)

	// set up the various commands to use
	BaseCli.AddCommand(keycmd.RootCmd)
	BaseCli.AddCommand(commands.InitCmd)
	BaseCli.AddCommand(seeds.RootCmd)
	proofs.StatePresenters.Register("account", AccountPresenter{})
	proofs.TxPresenters.Register("base", BaseTxPresenter{})
	BaseCli.AddCommand(proofs.RootCmd)
	txs.Register("send", SendTxMaker{})
	BaseCli.AddCommand(txs.RootCmd)
	BaseCli.AddCommand(proxy.RootCmd)
}

func main() {
	cmd := cli.PrepareMainCmd(BaseCli, "BC", os.ExpandEnv("$HOME/.basecli"))
	cmd.Execute()
}
