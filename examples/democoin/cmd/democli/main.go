package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/keys"
	"github.com/tepleton/tepleton-sdk/client/lcd"
	"github.com/tepleton/tepleton-sdk/client/rpc"
	"github.com/tepleton/tepleton-sdk/client/tx"

	"github.com/tepleton/tepleton-sdk/version"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/commands"
	bankcmd "github.com/tepleton/tepleton-sdk/x/bank/commands"
	ibccmd "github.com/tepleton/tepleton-sdk/x/ibc/commands"
	simplestakingcmd "github.com/tepleton/tepleton-sdk/x/simplestake/commands"

	"github.com/tepleton/tepleton-sdk/examples/democoin/app"
	"github.com/tepleton/tepleton-sdk/examples/democoin/types"
)

// toncliCmd is the entry point for this binary
var (
	democliCmd = &cobra.Command{
		Use:   "democli",
		Short: "Democoin light-client",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc, and tx commands
	rpc.AddCommands(democliCmd)
	democliCmd.AddCommand(client.LineBreak)
	tx.AddCommands(democliCmd, cdc)
	democliCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	democliCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetAccountDecoder(cdc)),
		)...)
	democliCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	democliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
		)...)
	democliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCRelayCmd(cdc),
			simplestakingcmd.BondTxCmd(cdc),
		)...)
	democliCmd.AddCommand(
		client.PostCommands(
			simplestakingcmd.UnbondTxCmd(cdc),
		)...)

	// add proxy, version and key info
	democliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(democliCmd, "BC", os.ExpandEnv("$HOME/.democli"))
	executor.Execute()
}
