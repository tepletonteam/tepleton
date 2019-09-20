package main

import (
	"os"

	"github.com/spf13/cobra"

	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/tepleton-sdk/client/commands"
	"github.com/tepleton/tepleton-sdk/client/commands/auto"
	"github.com/tepleton/tepleton-sdk/client/commands/proxy"
	"github.com/tepleton/tepleton-sdk/client/commands/query"
	rpccmd "github.com/tepleton/tepleton-sdk/client/commands/rpc"
	"github.com/tepleton/tepleton-sdk/client/commands/seeds"
	txcmd "github.com/tepleton/tepleton-sdk/client/commands/txs"
	eyescmd "github.com/tepleton/tepleton-sdk/modules/eyes/commands"
)

// EyesCli - main basecoin client command
var EyesCli = &cobra.Command{
	Use:   "eyescli",
	Short: "Light client for Tendermint",
	Long:  `EyesCli is the light client for a merkle key-value store (eyes)`,
}

func main() {
	commands.AddBasicFlags(EyesCli)

	// Prepare queries
	query.RootCmd.AddCommand(
		// These are default parsers, but optional in your app (you can remove key)
		query.TxQueryCmd,
		query.KeyQueryCmd,
		// this is our custom parser
		eyescmd.EyesQueryCmd,
	)

	// no middleware wrapers
	txcmd.Middleware = txcmd.Wrappers{}
	// txcmd.Middleware.Register(txcmd.RootCmd.PersistentFlags())

	// just the etc commands
	txcmd.RootCmd.AddCommand(
		eyescmd.SetTxCmd,
		eyescmd.RemoveTxCmd,
	)

	// Set up the various commands to use
	EyesCli.AddCommand(
		// we use out own init command to not require address arg
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		rpccmd.RootCmd,
		query.RootCmd,
		txcmd.RootCmd,
		proxy.RootCmd,
		commands.VersionCmd,
		auto.AutoCompleteCmd,
	)

	cmd := cli.PrepareMainCmd(EyesCli, "EYE", os.ExpandEnv("$HOME/.eyescli"))
	cmd.Execute()
}
