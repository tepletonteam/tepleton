package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tepleton/go-wire/data"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/types"
	"github.com/tepleton/tmlibs/log"
)

// ShowValidator - ported from Tendermint, show this node's validator info
func ShowValidatorCmd(logger log.Logger) *cobra.Command {
	cmd := showValidator{logger}
	return &cobra.Command{
		Use:   "show_validator",
		Short: "Show this node's validator info",
		RunE:  cmd.run,
	}
}

type showValidator struct {
	logger log.Logger
}

func (s showValidator) run(cmd *cobra.Command, args []string) error {
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}
	privValidator := types.LoadOrGenPrivValidatorFS(cfg.PrivValidatorFile())
	pubKeyJSONBytes, err := data.ToJSON(privValidator.PubKey)
	if err != nil {
		return err
	}
	fmt.Println(string(pubKeyJSONBytes))
	return nil
}
