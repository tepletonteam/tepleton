package commands

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	lc "github.com/tepleton/light-client"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/client/commands"
	proofcmd "github.com/tepleton/basecoin/client/commands/proofs"
	"github.com/tepleton/basecoin/modules/nonce"
	"github.com/tepleton/basecoin/stack"
)

// NonceQueryCmd - command to query an nonce account
var NonceQueryCmd = &cobra.Command{
	Use:   "nonce [address]",
	Short: "Get details of a nonce sequence number, with proof",
	RunE:  commands.RequireInit(nonceQueryCmd),
}

func nonceQueryCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing required argument [address]")
	}
	addr := strings.Join(args, ",")

	signers, err := commands.ParseActors(addr)
	if err != nil {
		return err
	}

	seq, proof, err := doNonceQuery(signers)
	if err != nil {
		return err
	}

	return proofcmd.OutputProof(seq, proof.BlockHeight())
}

func doNonceQuery(signers []basecoin.Actor) (sequence uint32, proof lc.Proof, err error) {

	key := stack.PrefixedKey(nonce.NameNonce, nonce.GetSeqKey(signers))

	proof, err = proofcmd.GetAndParseAppProof(key, &sequence)
	if lc.IsNoDataErr(err) {
		// no data, return sequence 0
		return 0, proof, nil
	}

	return
}
