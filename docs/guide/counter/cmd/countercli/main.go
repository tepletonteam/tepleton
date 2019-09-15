package main

import (
	"os"

	"github.com/spf13/cobra"

	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/basecoin/client/commands"
	"github.com/tepleton/basecoin/client/commands/proofs"
	"github.com/tepleton/basecoin/client/commands/proxy"
	"github.com/tepleton/basecoin/client/commands/seeds"
	txcmd "github.com/tepleton/basecoin/client/commands/txs"
	bcount "github.com/tepleton/basecoin/docs/guide/counter/cmd/countercli/commands"
	authcmd "github.com/tepleton/basecoin/modules/auth/commands"
	basecmd "github.com/tepleton/basecoin/modules/base/commands"
	coincmd "github.com/tepleton/basecoin/modules/coin/commands"
	feecmd "github.com/tepleton/basecoin/modules/fee/commands"
	noncecmd "github.com/tepleton/basecoin/modules/nonce/commands"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "countercli",
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
	proofs.RootCmd.AddCommand(
		// These are default parsers, optional in your app
		proofs.TxCmd,
		proofs.KeyCmd,
		coincmd.AccountQueryCmd,
		noncecmd.NonceQueryCmd,

		// XXX IMPORTANT: here is how you add custom query commands in your app
		bcount.CounterQueryCmd,
	)

	// set up the middleware
	txcmd.Middleware = txcmd.Wrappers{
		feecmd.FeeWrapper{},
		noncecmd.NonceWrapper{},
		basecmd.ChainWrapper{},
		authcmd.SigWrapper{},
	}
	txcmd.Middleware.Register(txcmd.RootCmd.PersistentFlags())

	// Prepare transactions
	proofs.TxPresenters.Register("base", txcmd.BaseTxPresenter{})
	txcmd.RootCmd.AddCommand(
		// This is the default transaction, optional in your app
		coincmd.SendTxCmd,

		// XXX IMPORTANT: here is how you add custom tx construction for your app
		bcount.CounterTxCmd,
	)

	// Set up the various commands to use
	BaseCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		proofs.RootCmd,
		txcmd.RootCmd,
		proxy.RootCmd,
	)

	cmd := cli.PrepareMainCmd(BaseCli, "CTL", os.ExpandEnv("$HOME/.countercli"))
	cmd.Execute()
}
