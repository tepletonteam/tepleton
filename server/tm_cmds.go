package server

import (
	"fmt"

	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/tepleton/tepleton-sdk/types"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/p2p"
	pvm "github.com/tepleton/tepleton/types/priv_validator"
)

// ShowNodeIDCmd - ported from Tendermint, dump node ID to stdout
func ShowNodeIDCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "show_node_id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidator - ported from Tendermint, show this node's validator info
func ShowValidatorCmd(ctx *Context) *cobra.Command {
	flagJSON := "json"
	cmd := cobra.Command{
		Use:   "show_validator",
		Short: "Show this node's tepleton validator info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := ctx.Config
			privValidator := pvm.LoadOrGenFilePV(cfg.PrivValidatorFile())
			valAddr := sdk.Address(privValidator.PubKey.Address())

			if viper.GetBool(flagJSON) {

				cdc := wire.NewCodec()
				wire.RegisterCrypto(cdc)
				pubKeyJSONBytes, err := cdc.MarshalJSON(valAddr)
				if err != nil {
					return err
				}
				fmt.Println(string(pubKeyJSONBytes))
				return nil
			}
			addr, err := sdk.Bech32TepletonifyVal(valAddr)
			if err != nil {
				return err
			}
			fmt.Println(addr)
			return nil
		},
	}
	cmd.Flags().Bool(flagJSON, false, "get machine parseable output")
	return &cmd
}

// UnsafeResetAllCmd - extension of the tepleton command, resets initialization
func UnsafeResetAllCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe_reset_all",
		Short: "Reset blockchain database, priv_validator.json file, and the logger",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			tcmd.ResetAll(cfg.DBDir(), cfg.PrivValidatorFile(), ctx.Logger)
			return nil
		},
	}
}
