package main

import (
	"github.com/spf13/cobra"

	"github.com/tepleton/tepleton/libs/cli"

	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/keys"
	"github.com/tepleton/tepleton-sdk/client/lcd"
	"github.com/tepleton/tepleton-sdk/client/rpc"
	"github.com/tepleton/tepleton-sdk/client/tx"
	"github.com/tepleton/tepleton-sdk/version"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/client/cli"
	bankcmd "github.com/tepleton/tepleton-sdk/x/bank/client/cli"
	govcmd "github.com/tepleton/tepleton-sdk/x/gov/client/cli"
	ibccmd "github.com/tepleton/tepleton-sdk/x/ibc/client/cli"
	slashingcmd "github.com/tepleton/tepleton-sdk/x/slashing/client/cli"
	stakecmd "github.com/tepleton/tepleton-sdk/x/stake/client/cli"

	"github.com/tepleton/tepleton-sdk/cmd/ton/app"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "toncli",
		Short: "Gaia light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc commands
	rpc.AddCommands(rootCmd)

	//Add state commands
	tepletonCmd := &cobra.Command{
		Use:   "tepleton",
		Short: "Tendermint state querying subcommands",
	}
	tepletonCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tepletonCmd, cdc)

	//Add IBC commands
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	ibcCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
		)...)

	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced subcommands",
	}

	advancedCmd.AddCommand(
		tepletonCmd,
		ibcCmd,
		lcd.ServeCommand(cdc),
	)
	rootCmd.AddCommand(
		advancedCmd,
		client.LineBreak,
	)

	//Add stake commands
	stakeCmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake and validation subcommands",
	}
	stakeCmd.AddCommand(
		client.GetCommands(
			stakecmd.GetCmdQueryValidator("stake", cdc),
			stakecmd.GetCmdQueryValidators("stake", cdc),
			stakecmd.GetCmdQueryDelegation("stake", cdc),
			stakecmd.GetCmdQueryDelegations("stake", cdc),
			slashingcmd.GetCmdQuerySigningInfo("slashing", cdc),
		)...)
	stakeCmd.AddCommand(
		client.PostCommands(
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond("stake", cdc),
			stakecmd.GetCmdRedelegate("stake", cdc),
			slashingcmd.GetCmdUnrevoke(cdc),
		)...)
	rootCmd.AddCommand(
		stakeCmd,
	)

	//Add stake commands
	govCmd := &cobra.Command{
		Use:   "gov",
		Short: "Governance and voting subcommands",
	}
	govCmd.AddCommand(
		client.GetCommands(
			govcmd.GetCmdQueryProposal("gov", cdc),
			govcmd.GetCmdQueryVote("gov", cdc),
		)...)
	govCmd.AddCommand(
		client.PostCommands(
			govcmd.GetCmdSubmitProposal(cdc),
			govcmd.GetCmdDeposit(cdc),
			govcmd.GetCmdVote(cdc),
		)...)
	rootCmd.AddCommand(
		govCmd,
	)

	//Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
