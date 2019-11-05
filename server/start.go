package server

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/wrsp/server"

	"github.com/tepleton/tepleton-sdk/baseapp"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	pvm "github.com/tepleton/tepleton/types/priv_validator"
	cmn "github.com/tepleton/tmlibs/common"
)

const (
	flagWithTendermint = "with-tepleton"
	flagAddress        = "address"
)

// StartCmd runs the service passed in, either
// stand-alone, or in-process with tepleton
func StartCmd(ctx *Context, appCreator baseapp.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.GetBool(flagWithTendermint) {
				ctx.Logger.Info("Starting WRSP without Tendermint")
				return startStandAlone(ctx, appCreator)
			}
			ctx.Logger.Info("Starting WRSP with Tendermint")
			return startInProcess(ctx, appCreator)
		},
	}

	// basic flags for wrsp app
	cmd.Flags().Bool(flagWithTendermint, true, "run wrsp app embedded in-process with tepleton")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:46658", "Listen address")

	// AddNodeFlags adds support for all tepleton-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func startStandAlone(ctx *Context, appCreator baseapp.AppCreator) error {
	// Generate the app in the proper dir
	addr := viper.GetString(flagAddress)
	home := viper.GetString("home")
	app, err := appCreator(home, ctx.Logger)
	if err != nil {
		return err
	}

	svr, err := server.NewServer(addr, "socket", app)
	if err != nil {
		return errors.Errorf("Error creating listener: %v\n", err)
	}
	svr.SetLogger(ctx.Logger.With("module", "wrsp-server"))
	svr.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}

func startInProcess(ctx *Context, appCreator baseapp.AppCreator) error {
	cfg := ctx.Config
	home := cfg.RootDir
	app, err := appCreator(home, ctx.Logger)
	if err != nil {
		return err
	}

	// Create & start tepleton node
	n, err := node.NewNode(cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		ctx.Logger.With("module", "node"))
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
