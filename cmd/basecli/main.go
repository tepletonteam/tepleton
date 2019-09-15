package main

import (
	"os"

	"github.com/spf13/cobra"

	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/basecoin/client/commands"
	"github.com/tepleton/basecoin/client/commands/auto"
	"github.com/tepleton/basecoin/client/commands/proofs"
	"github.com/tepleton/basecoin/client/commands/proxy"
	rpccmd "github.com/tepleton/basecoin/client/commands/rpc"
	"github.com/tepleton/basecoin/client/commands/seeds"
	txcmd "github.com/tepleton/basecoin/client/commands/txs"
	authcmd "github.com/tepleton/basecoin/modules/auth/commands"
	basecmd "github.com/tepleton/basecoin/modules/base/commands"
	coincmd "github.com/tepleton/basecoin/modules/coin/commands"
	feecmd "github.com/tepleton/basecoin/modules/fee/commands"
	noncecmd "github.com/tepleton/basecoin/modules/nonce/commands"
	rolecmd "github.com/tepleton/basecoin/modules/roles/commands"
)

// BaseCli - main basecoin client command
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
	proofs.RootCmd.AddCommand(
		// These are default parsers, but optional in your app (you can remove key)
		proofs.TxQueryCmd,
		proofs.KeyQueryCmd,
		coincmd.AccountQueryCmd,
		noncecmd.NonceQueryCmd,
		rolecmd.RoleQueryCmd,
	)
	proofs.TxPresenters.Register("base", txcmd.BaseTxPresenter{})

	// set up the middleware
	txcmd.Middleware = txcmd.Wrappers{
		feecmd.FeeWrapper{},
		rolecmd.RoleWrapper{},
		noncecmd.NonceWrapper{},
		basecmd.ChainWrapper{},
		authcmd.SigWrapper{},
	}
	txcmd.Middleware.Register(txcmd.RootCmd.PersistentFlags())

	// you will always want this for the base send command
	txcmd.RootCmd.AddCommand(
		// This is the default transaction, optional in your app
		coincmd.SendTxCmd,
		// this enables creating roles
		rolecmd.CreateRoleTxCmd,
	)

	// Set up the various commands to use
	BaseCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		rpccmd.RootCmd,
		proofs.RootCmd,
		txcmd.RootCmd,
		proxy.RootCmd,
		commands.VersionCmd,
		auto.AutoCompleteCmd,
	)

	cmd := cli.PrepareMainCmd(BaseCli, "BC", os.ExpandEnv("$HOME/.basecli"))
	cmd.Execute()
}
