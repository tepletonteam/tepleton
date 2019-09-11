package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/wrsp/server"
	"github.com/tepleton/basecoin"
	eyesApp "github.com/tepleton/merkleeyes/app"
	eyes "github.com/tepleton/merkleeyes/client"
	"github.com/tepleton/tmlibs/cli"
	cmn "github.com/tepleton/tmlibs/common"

	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	"github.com/tepleton/tepleton/types"

	"github.com/tepleton/basecoin/app"
	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/stack"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start basecoin",
	RunE:  startCmd,
}

// TODO: move to config file
const EyesCacheSize = 10000

//nolint
const (
	FlagAddress           = "address"
	FlagEyes              = "eyes"
	FlagWithoutTendermint = "without-tepleton"
)

func init() {
	flags := StartCmd.Flags()
	flags.String(FlagAddress, "tcp://0.0.0.0:46658", "Listen address")
	flags.String(FlagEyes, "local", "MerkleEyes address, or 'local' for embedded")
	flags.Bool(FlagWithoutTendermint, false, "Only run basecoin wrsp app, assume external tepleton process")
	// add all standard 'tepleton node' flags
	tcmd.AddNodeFlags(StartCmd)
}

// TODO: setup handler instead of Plugins
func getHandler() basecoin.Handler {
	// use the default stack
	h := coin.NewHandler()
	app := stack.NewDefault().Use(h)
	return app

	// register IBC plugn
	// basecoinApp.RegisterPlugin(NewIBCPlugin())

	// register all other plugins
	// for _, p := range plugins {
	//  basecoinApp.RegisterPlugin(p.newPlugin())
	// }
}

func startCmd(cmd *cobra.Command, args []string) error {
	rootDir := viper.GetString(cli.HomeFlag)
	meyes := viper.GetString(FlagEyes)

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if meyes == "local" {
		eyesApp.SetLogger(logger.With("module", "merkleeyes"))
		eyesCli = eyes.NewLocalClient(path.Join(rootDir, "data", "merkleeyes.db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(meyes)
		if err != nil {
			return errors.Errorf("Error connecting to MerkleEyes: %v\n", err)
		}
	}

	// Create Basecoin app
	h := app.DefaultHandler()
	basecoinApp := app.NewBasecoin(h, eyesCli, logger.With("module", "app"))

	// if chain_id has not been set yet, load the genesis.
	// else, assume it's been loaded
	if basecoinApp.GetState().GetChainID() == "" {
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

	chainID := basecoinApp.GetState().GetChainID()
	if viper.GetBool(FlagWithoutTendermint) {
		logger.Info("Starting Basecoin without Tendermint", "chain_id", chainID)
		// run just the wrsp app/server
		return startBasecoinWRSP(basecoinApp)
	} else {
		logger.Info("Starting Basecoin with Tendermint", "chain_id", chainID)
		// start the app with tepleton in-process
		return startTendermint(rootDir, basecoinApp)
	}
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
	privValidator := types.LoadOrGenPrivValidator(cfg.PrivValidatorFile(), logger)
	n := node.NewNode(cfg, privValidator, proxy.NewLocalClientCreator(basecoinApp), logger.With("module", "node"))

	_, err = n.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
