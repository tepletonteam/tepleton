package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/basecoin/client/commands"
	proofcmd "github.com/tepleton/basecoin/client/commands/proofs"
	"github.com/tepleton/basecoin/modules/roles"
	"github.com/tepleton/basecoin/stack"
)

// RoleQueryCmd - command to query a role
var RoleQueryCmd = &cobra.Command{
	Use:   "role [name]",
	Short: "Get details of a role, with proof",
	RunE:  commands.RequireInit(roleQueryCmd),
}

func roleQueryCmd(cmd *cobra.Command, args []string) error {
	arg, err := commands.GetOneArg(args, "name")
	if err != nil {
		return err
	}
	role, err := parseRole(arg)
	if err != nil {
		return err
	}

	var res roles.Role
	key := stack.PrefixedKey(roles.NameRole, role)
	prove := !viper.GetBool(commands.FlagTrustNode)
	height, err := proofcmd.GetParsed(key, &res, prove)
	if err != nil {
		return err
	}

	return proofcmd.OutputProof(res, height)
}
