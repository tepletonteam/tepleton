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

	bcmd "github.com/tepleton/basecoin/cmd/basecli/commands"
	bcount "github.com/tepleton/basecoin/cmd/countercli/commands"
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

	// Prepare queries
	pr := proofs.RootCmd
	// These are default parsers, but you optional in your app
	pr.AddCommand(proofs.TxCmd)
	pr.AddCommand(proofs.KeyCmd)
	pr.AddCommand(bcmd.AccountQueryCmd)

	// IMPORTANT: here is how you add custom query commands in your app
	pr.AddCommand(bcount.CounterQueryCmd)

	proofs.TxPresenters.Register("base", bcmd.BaseTxPresenter{})
	tr := txs.RootCmd
	tr.AddCommand(bcmd.SendTxCmd)

	// IMPORTANT: here is how you add custom tx construction for your app
	tr.AddCommand(bcount.CounterTxCmd)

	// Set up the various commands to use
	BaseCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		pr,
		tr,
		proxy.RootCmd)

	cmd := cli.PrepareMainCmd(BaseCli, "BC", os.ExpandEnv("$HOME/.basecli"))
	cmd.Execute()
}
