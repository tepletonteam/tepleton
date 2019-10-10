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
	coolcmd "github.com/tepleton/tepleton-sdk/examples/basecoin/x/cool/commands"
	"github.com/tepleton/tepleton-sdk/version"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/commands"
	bankcmd "github.com/tepleton/tepleton-sdk/x/bank/commands"

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
			coolcmd.WhatCoolTxCmd(cdc),
		)...)
	basecliCmd.AddCommand(
		client.PostCommands(
			coolcmd.SetWhatCoolTxCmd(cdc),
		)...)

	// add proxy, version and key info
	basecliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(basecliCmd, "BC", os.ExpandEnv("$HOME/.basecli"))
	executor.Execute()
}
