package server

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/wrsp/server"

	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	"github.com/tepleton/tepleton/types"
	cmn "github.com/tepleton/tmlibs/common"

	"github.com/tepleton/tepleton-sdk/baseapp"
)

const (
	flagWithTendermint = "with-tepleton"
	flagAddress        = "address"
)

// StartCmd runs the service passed in, either
// stand-alone, or in-process with tepleton
func StartCmd(app *baseapp.BaseApp) *cobra.Command {
	start := startCmd{
		app: app,
	}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		RunE:  start.run,
	}
	// basic flags for wrsp app
	cmd.Flags().Bool(flagWithTendermint, true, "run wrsp app embedded in-process with tepleton")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:46658", "Listen address")

	// AddNodeFlags adds support for all
	// tepleton-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

type startCmd struct {
	// do this in main:
	// rootDir := viper.GetString(cli.HomeFlag)
	// node.Logger = ....
	app *baseapp.BaseApp
}

func (s startCmd) run(cmd *cobra.Command, args []string) error {
	logger := s.app.Logger
	if !viper.GetBool(flagWithTendermint) {
		logger.Info("Starting WRSP without Tendermint")
		return s.startStandAlone()
	}
	logger.Info("Starting WRSP with Tendermint")
	return s.startInProcess()
}

func (s startCmd) startStandAlone() error {
	logger := s.app.Logger
	// Start the WRSP listener
	addr := viper.GetString(flagAddress)
	svr, err := server.NewServer(addr, "socket", s.app)
	if err != nil {
		return errors.Errorf("Error creating listener: %v\n", err)
	}
	svr.SetLogger(logger.With("module", "wrsp-server"))
	svr.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}

func (s startCmd) startInProcess() error {
	logger := s.app.Logger
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}

	// Create & start tepleton node
	n, err := node.NewNode(cfg,
		types.LoadOrGenPrivValidatorFS(cfg.PrivValidatorFile()),
		proxy.NewLocalClientCreator(s.app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		logger.With("module", "node"))
	if err != nil {
		return err
	}

	err = n.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
