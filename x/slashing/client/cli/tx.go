package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/client/cli"
	"github.com/tepleton/tepleton-sdk/x/slashing"
)

// create unrevoke command
func GetCmdUnrevoke(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unrevoke",
		Args:  cobra.ExactArgs(1),
		Short: "unrevoke validator previously revoked for downtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			validatorAddr, err := sdk.GetAccAddressBech32(args[0])
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnrevoke(validatorAddr)

			// build and sign the transaction, then broadcast to Tendermint
			res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, cdc)
			if err != nil {
				return err
			}

			fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return nil
		},
	}
	return cmd
}
