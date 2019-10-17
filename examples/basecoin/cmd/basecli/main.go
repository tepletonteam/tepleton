package main

import (
	"errors"
	"github.com/spf13/cobra"
	"os"

	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/keys"
	"github.com/tepleton/tepleton-sdk/client/lcd"
	"github.com/tepleton/tepleton-sdk/client/rpc"
	"github.com/tepleton/tepleton-sdk/client/tx"

	coolcmd "github.com/tepleton/tepleton-sdk/examples/basecoin/x/cool/commands"
	"github.com/tepleton/tepleton-sdk/version"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/commands"
	bankcmd "github.com/tepleton/tepleton-sdk/x/bank/commands"
	ibccmd "github.com/tepleton/tepleton-sdk/x/ibc/commands"
	stakingcmd "github.com/tepleton/tepleton-sdk/x/staking/commands"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	"github.com/tepleton/tepleton-sdk/examples/basecoin/types"
)

// toncliCmd is the entry point for this binary
var (
	basecliCmd = &cobra.Command{
		Use:   "basecli",
		Short: "Basecoin light-client",
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
	rpc.AddCommands(basecliCmd)
	basecliCmd.AddCommand(client.LineBreak)
	tx.AddCommands(basecliCmd, cdc)
	basecliCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	basecliCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetParseAccount(cdc)),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			coolcmd.QuizTxCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			coolcmd.SetTrendTxCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCRelayCmd(cdc),
			stakingcmd.BondTxCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			stakingcmd.UnbondTxCmd(cdc),
		)...)

	// add proxy, version and key info
	basecliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(basecliCmd, "BC", os.ExpandEnv("$HOME/.basecli"))
	executor.Execute()
}
