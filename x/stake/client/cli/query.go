package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	crypto "github.com/tepleton/go-crypto"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire" // XXX fix
	"github.com/tepleton/tepleton-sdk/x/stake"
)

// get the command to query a validator
func GetCmdQueryValidator(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "Query a validator-validator account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(args[0])
			if err != nil {
				return err
			}
			key := stake.GetValidatorKey(addr)
			ctx := context.NewCoreContextFromViper()
			res, err := ctx.Query(key, storeName)
			if err != nil {
				return err
			}

			// parse out the validator
			validator := new(stake.Validator)
			cdc.MustUnmarshalBinary(res, validator)
			err = cdc.UnmarshalBinary(res, validator)
			if err != nil {
				return err
			}
			output, err := wire.MarshalJSONIndent(cdc, validator)
			fmt.Println(string(output))

			// TODO output with proofs / machine parseable etc.
		},
	}

	return cmd
}

// get the command to query a candidate
func GetCmdQueryCandidates(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "candidates",
		Short: "Query for all validator-candidate accounts",
		RunE: func(cmd *cobra.Command, args []string) error {

			key := stake.CandidatesKey
			ctx := context.NewCoreContextFromViper()
			resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
			if err != nil {
				return err
			}

			// parse out the candidates
			var candidates []stake.Candidate
			for _, KV := range resKVs {
				var candidate stake.Candidate
				cdc.MustUnmarshalBinary(KV.Value, &candidate)
				candidates = append(candidates, candidate)
			}

			output, err := wire.MarshalJSONIndent(cdc, candidates)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	return cmd
}

// get the command to query a single delegation bond
func GetCmdQueryDelegation(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-bond",
		Short: "Query a delegations bond based on address and validator pubkey",
		RunE: func(cmd *cobra.Command, args []string) error {

			addr, err := sdk.GetAddress(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}

			bz, err := hex.DecodeString(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}
			delegation := crypto.Address(bz)

			key := stake.GetDelegationKey(delegation, addr, cdc)
			ctx := context.NewCoreContextFromViper()
			res, err := ctx.Query(key, storeName)
			if err != nil {
				return err
			}

			// parse out the bond
			bond := new(stake.Delegation)
			cdc.MustUnmarshalBinary(res, bond)
			output, err := wire.MarshalJSONIndent(cdc, bond)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)
	return cmd
}

// get the command to query all the candidates bonded to a delegation
func GetCmdQueryDelegations(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-candidates",
		Short: "Query all delegations bonds based on delegation-address",
		RunE: func(cmd *cobra.Command, args []string) error {

			delegatorAddr, err := sdk.GetAddress(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}
			key := stake.GetDelegatorBondsKey(delegatorAddr, cdc)
			ctx := context.NewCoreContextFromViper()
			resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
			if err != nil {
				return err
			}

			// parse out the candidates
			var delegations []stake.Delegation
			for _, KV := range resKVs {
				var delegation stake.Delegation
				cdc.MustUnmarshalBinary(KV.Value, &delegation)
				delegations = append(delegations, delegation)
			}

			output, err := wire.MarshalJSONIndent(cdc, delegations)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil

			// TODO output with proofs / machine parseable etc.
		},
	}
	cmd.Flags().AddFlagSet(fsDelegator)
	return cmd
}
