package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	txcmd "github.com/tepleton/light-client/commands/txs"

	"github.com/tepleton/basecoin"
	bcmd "github.com/tepleton/basecoin/cmd/basecli/commands"
	"github.com/tepleton/basecoin/docs/guide/counter/plugins/counter"
	"github.com/tepleton/basecoin/modules/auth"
	"github.com/tepleton/basecoin/modules/coin"
)

//CounterTxCmd is the CLI command to execute the counter
//  through the appTx Command
var CounterTxCmd = &cobra.Command{
	Use:   "counter",
	Short: "add a vote to the counter",
	Long: `Add a vote to the counter.

You must pass --valid for it to count and the countfee will be added to the counter.`,
	RunE: counterTx,
}

// nolint - flags names
const (
	FlagCountFee = "countfee"
	FlagValid    = "valid"
)

func init() {
	fs := CounterTxCmd.Flags()
	fs.String(FlagCountFee, "", "Coins to send in the format <amt><coin>,<amt><coin>...")
	fs.Bool(FlagValid, false, "Is count valid?")

	fs.String(bcmd.FlagFee, "0mycoin", "Coins for the transaction fee of the format <amt><coin>")
	fs.Int(bcmd.FlagSequence, -1, "Sequence number for this transaction")
}

// TODO: counterTx is very similar to the sendtx one,
// maybe we can pull out some common patterns?
func counterTx(cmd *cobra.Command, args []string) error {
	// load data from json or flags
	var tx basecoin.Tx
	found, err := txcmd.LoadJSON(&tx)
	if err != nil {
		return err
	}
	if !found {
		tx, err = readCounterTxFlags()
	}
	if err != nil {
		return err
	}

	// TODO: make this more flexible for middleware
	tx, err = bcmd.WrapFeeTx(tx)
	if err != nil {
		return err
	}
	tx, err = bcmd.WrapChainTx(tx)
	if err != nil {
		return err
	}
	tx, err = bcmd.WrapNonceTx(tx)
	if err != nil {
		return err
	}

	stx := auth.NewSig(tx)

	// Sign if needed and post.  This it the work-horse
	bres, err := txcmd.SignAndPostTx(stx)
	if err != nil {
		return err
	}

	// Output result
	return txcmd.OutputTx(bres)
}

func readCounterTxFlags() (tx basecoin.Tx, err error) {
	feeCoins, err := coin.ParseCoins(viper.GetString(FlagCountFee))
	if err != nil {
		return tx, err
	}

	tx = counter.NewTx(viper.GetBool(FlagValid), feeCoins)
	return tx, nil
}
