package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tepleton/basecoin"
	wire "github.com/tepleton/go-wire"
	lc "github.com/tepleton/light-client"
	lcmd "github.com/tepleton/light-client/commands"
	proofcmd "github.com/tepleton/light-client/commands/proofs"
	"github.com/tepleton/light-client/proofs"

	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/stack"
)

// AccountQueryCmd - command to query an account
var AccountQueryCmd = &cobra.Command{
	Use:   "account [address]",
	Short: "Get details of an account, with proof",
	RunE:  lcmd.RequireInit(doAccountQuery),
}

func doAccountQuery(cmd *cobra.Command, args []string) error {
	addr, err := proofcmd.ParseHexKey(args, "address")
	if err != nil {
		return err
	}
	key := coin.NewAccountant("").MakeKey(stack.SigPerm(addr))

	acc := coin.Account{}
	proof, err := proofcmd.GetAndParseAppProof(key, &acc)
	if lc.IsNoDataErr(err) {
		return errors.Errorf("Account bytes are empty for address %X ", addr)
	} else if err != nil {
		return err
	}

	return proofcmd.OutputProof(acc, proof.BlockHeight())
}

// BaseTxPresenter this decodes all basecoin tx
type BaseTxPresenter struct {
	proofs.RawPresenter // this handles MakeKey as hex bytes
}

// ParseData - parse BaseTxPresenter Data
func (b BaseTxPresenter) ParseData(raw []byte) (interface{}, error) {
	var tx basecoin.Tx
	err := wire.ReadBinaryBytes(raw, &tx)
	return tx, err
}
