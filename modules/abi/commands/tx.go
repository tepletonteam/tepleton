package commands

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/basecoin/client/commands"
	txcmd "github.com/tepleton/basecoin/client/commands/txs"
	"github.com/tepleton/basecoin/modules/abi"
	"github.com/tepleton/light-client/certifiers"
)

// RegisterChainTxCmd is CLI command to register a new chain for abi
var RegisterChainTxCmd = &cobra.Command{
	Use:   "abi-register",
	Short: "Register a new chain",
	RunE:  commands.RequireInit(registerChainTxCmd),
}

// UpdateChainTxCmd is CLI command to update the header for an abi chain
var UpdateChainTxCmd = &cobra.Command{
	Use:   "abi-update",
	Short: "Add new header to an existing chain",
	RunE:  commands.RequireInit(updateChainTxCmd),
}

// PostPacketTxCmd is CLI command to post abi packet on the destination chain
var PostPacketTxCmd = &cobra.Command{
	Use:   "abi-post",
	Short: "Post an abi packet on the destination chain",
	RunE:  commands.RequireInit(postPacketTxCmd),
}

// TODO: relay!

//nolint
const (
	FlagSeed   = "seed"
	FlagPacket = "packet"
)

func init() {
	fs1 := RegisterChainTxCmd.Flags()
	fs1.String(FlagSeed, "", "Filename with a seed file")

	fs2 := UpdateChainTxCmd.Flags()
	fs2.String(FlagSeed, "", "Filename with a seed file")

	fs3 := PostPacketTxCmd.Flags()
	fs3.String(FlagPacket, "", "Filename with a packet to post")
}

func registerChainTxCmd(cmd *cobra.Command, args []string) error {
	seed, err := readSeed()
	if err != nil {
		return err
	}
	tx := abi.RegisterChainTx{seed}.Wrap()
	return txcmd.DoTx(tx)
}

func updateChainTxCmd(cmd *cobra.Command, args []string) error {
	seed, err := readSeed()
	if err != nil {
		return err
	}
	tx := abi.UpdateChainTx{seed}.Wrap()
	return txcmd.DoTx(tx)
}

func postPacketTxCmd(cmd *cobra.Command, args []string) error {
	post, err := readPostPacket()
	if err != nil {
		return err
	}
	return txcmd.DoTx(post.Wrap())
}

func readSeed() (seed certifiers.Seed, err error) {
	name := viper.GetString(FlagSeed)
	if name == "" {
		return seed, errors.New("You must specify a seed file")
	}

	err = readFile(name, &seed)
	return
}

func readPostPacket() (post abi.PostPacketTx, err error) {
	name := viper.GetString(FlagPacket)
	if name == "" {
		return post, errors.New("You must specify a packet file")
	}

	err = readFile(name, &post)
	return
}

func readFile(name string, input interface{}) (err error) {
	var f *os.File
	f, err = os.Open(name)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	// read the file as json into a seed
	j := json.NewDecoder(f)
	err = j.Decode(input)
	return errors.Wrap(err, "Invalid file")
}
