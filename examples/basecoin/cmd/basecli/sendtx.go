package main

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/bank"
	crypto "github.com/tepleton/go-crypto"
	"github.com/tepleton/tmlibs/cli"
)

const (
	flagTo       = "to"
	flagAmount   = "amount"
	flagFee      = "fee"
	flagSequence = "seq"
)

func postSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE:  sendTx,
	}
	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagFee, "", "Fee to pay along with transaction")
	cmd.Flags().Int64(flagSequence, 0, "Sequence number to sign the tx")
	return cmd
}

func sendTx(cmd *cobra.Command, args []string) error {
	txBytes, err := buildTx()
	if err != nil {
		return err
	}

	uri := viper.GetString(flagNode)
	if uri == "" {
		return errors.New("Must define which node to query with --node")
	}
	node := client.GetNode(uri)

	result, err := node.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}

	if result.CheckTx.Code != uint32(0) {
		return errors.Errorf("CheckTx failed: (%d) %s",
			result.CheckTx.Code,
			result.CheckTx.Log)
	}
	if result.DeliverTx.Code != uint32(0) {
		return errors.Errorf("CheckTx failed: (%d) %s",
			result.DeliverTx.Code,
			result.DeliverTx.Log)
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", result.Height, result.Hash.String())
	return nil
}

func buildTx() ([]byte, error) {
	rootDir := viper.GetString(cli.HomeFlag)
	keybase, err := client.GetKeyBase(rootDir)
	if err != nil {
		return nil, err
	}

	name := viper.GetString(flagName)
	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.WithMessage(err, "No key for name")
	}
	from := info.PubKey.Address()

	msg, err := buildMsg(from)
	if err != nil {
		return nil, err
	}

	// sign and build
	bz := msg.GetSignBytes()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	passphrase, err := client.GetPassword(prompt)
	if err != nil {
		return nil, err
	}
	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  viper.GetInt64(flagSequence),
	}}

	// marshal bytes
	tx := sdk.NewStdTx(msg, sigs)
	cdc := app.MakeTxCodec()
	txBytes, err := cdc.MarshalBinary(tx)
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

func buildMsg(from crypto.Address) (sdk.Msg, error) {
	// parse coins
	amount := viper.GetString(flagAmount)
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, err
	}

	// parse destination address
	dest := viper.GetString(flagTo)
	bz, err := hex.DecodeString(dest)
	if err != nil {
		return nil, err
	}
	to := crypto.Address(bz)

	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewSendMsg([]bank.Input{input}, []bank.Output{output})
	return msg, nil
}
