package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/client/cli"

	"github.com/tepleton/tepleton-sdk/examples/democoin/x/cool"
)

// take the coolness quiz transaction
func QuizTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cool [answer]",
		Short: "What's cooler than being cool?",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from address from the name flag
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			// create the message
			msg := cool.NewMsgQuiz(from, args[0])

			// get account name
			name := viper.GetString(client.FlagName)

			// build and sign the transaction, then broadcast to Tendermint
			err = ctx.EnsureSignBuildBroadcast(name, []sdk.Msg{msg}, cdc, ctx.Async, false)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

// set a new cool trend transaction
func SetTrendTxCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setcool [answer]",
		Short: "You're so cool, tell us what is cool!",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from address from the name flag
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			// get account name
			name := viper.GetString(client.FlagName)

			// create the message
			msg := cool.NewMsgSetTrend(from, args[0])

			// build and sign the transaction, then broadcast to Tendermint
			err = ctx.EnsureSignBuildBroadcast(name, []sdk.Msg{msg}, cdc, ctx.Async, false)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
