package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/wrsp/server"
	"github.com/tepleton/tmlibs/cli"
	cmn "github.com/tepleton/tmlibs/common"

	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	"github.com/tepleton/tepleton/types"

	"github.com/tepleton/tepleton-sdk/app"
)

// StartCmd - command to start running the wrsp app (and tepleton)!
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start this full node",
	RunE:  startCmd,
}

// nolint TODO: move to config file
const EyesCacheSize = 10000

//nolint
const (
	FlagAddress           = "address"
	FlagWithoutTendermint = "without-tepleton"
)

var (
	// Handler - use a global to store the handler, so we can set it in main.
	// TODO: figure out a cleaner way to register plugins
	Handler sdk.Handler
)

func init() {
	flags := StartCmd.Flags()
	flags.String(FlagAddress, "tcp://0.0.0.0:46658", "Listen address")
	flags.Bool(FlagWithoutTendermint, false, "Only run wrsp app, assume external tepleton process")
	// add all standard 'tepleton node' flags
	tcmd.AddNodeFlags(StartCmd)
}

func startCmd(cmd *cobra.Command, args []string) error {
	rootDir := viper.GetString(cli.HomeFlag)

	store, err := app.NewStore(
		path.Join(rootDir, "data", "merkleeyes.db"),
		EyesCacheSize,
		logger.With("module", "store"),
	)
	if err != nil {
		return err
	}

	// Create Basecoin app
	basecoinApp := app.NewBasecoin(Handler, store, logger.With("module", "app"))

	// if chain_id has not been set yet, load the genesis.
	// else, assume it's been loaded
	if basecoinApp.GetChainID() == "" {
		// If genesis file exists, set key-value options
		genesisFile := path.Join(rootDir, "genesis.json")
		if _, err := os.Stat(genesisFile); err == nil {
			err := basecoinApp.LoadGenesis(genesisFile)
			if err != nil {
				return errors.Errorf("Error in LoadGenesis: %v\n", err)
			}
		} else {
			fmt.Printf("No genesis file at %s, skipping...\n", genesisFile)
		}
	}

	chainID := basecoinApp.GetChainID()
	if viper.GetBool(FlagWithoutTendermint) {
		logger.Info("Starting Basecoin without Tendermint", "chain_id", chainID)
		// run just the wrsp app/server
		return startBasecoinWRSP(basecoinApp)
	}
	logger.Info("Starting Basecoin with Tendermint", "chain_id", chainID)
	// start the app with tepleton in-process
	return startTendermint(rootDir, basecoinApp)
}

func startBasecoinWRSP(basecoinApp *app.Basecoin) error {
	// Start the WRSP listener
	addr := viper.GetString(FlagAddress)
	svr, err := server.NewServer(addr, "socket", basecoinApp)
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

func startTendermint(dir string, basecoinApp *app.Basecoin) error {
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}

	// Create & start tepleton node
	n, err := node.NewNode(cfg,
		types.LoadOrGenPrivValidatorFS(cfg.PrivValidatorFile()),
		proxy.NewLocalClientCreator(basecoinApp),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		logger.With("module", "node"))
	if err != nil {
		return err
	}

	_, err = n.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
