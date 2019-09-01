package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/wrsp/server"
	eyes "github.com/tepleton/merkleeyes/client"
	"github.com/tepleton/tmlibs/cli"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton/config"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	"github.com/tepleton/tepleton/types"

	"github.com/tepleton/basecoin/app"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start basecoin",
	RunE:  startCmd,
}

//flags
var (
	addrFlag              string
	eyesFlag              string
	dirFlag               string
	withoutTendermintFlag bool
)

// TODO: move to config file
const EyesCacheSize = 10000

func init() {

	flags := []Flag2Register{
		{&addrFlag, "address", "tcp://0.0.0.0:46658", "Listen address"},
		{&eyesFlag, "eyes", "local", "MerkleEyes address, or 'local' for embedded"},
		{&dirFlag, "dir", ".", "Root directory"},
		{&withoutTendermintFlag, "without-tepleton", false, "Run Tendermint in-process with the App"},
	}
	RegisterFlags(StartCmd, flags)
}

func startCmd(cmd *cobra.Command, args []string) error {
	rootDir := viper.GetString(cli.HomeFlag)

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if eyesFlag == "local" {
		eyesCli = eyes.NewLocalClient(path.Join(rootDir, "data", "merkleeyes.db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(eyesFlag)
		if err != nil {
			return errors.Errorf("Error connecting to MerkleEyes: %v\n", err)
		}
	}

	// Create Basecoin app
	basecoinApp := app.NewBasecoin(eyesCli)
	basecoinApp.SetLogger(logger.With("module", "app"))

	// register IBC plugn
	basecoinApp.RegisterPlugin(NewIBCPlugin())

	// register all other plugins
	for _, p := range plugins {
		basecoinApp.RegisterPlugin(p.newPlugin())
	}

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
	if withoutTendermintFlag {
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
	svr, err := server.NewServer(addrFlag, "socket", basecoinApp)
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
	cfg := config.DefaultConfig()
	err := viper.Unmarshal(cfg)
	if err != nil {
		return err
	}
	cfg.SetRoot(cfg.RootDir)
	config.EnsureRoot(cfg.RootDir)

	tmLogger, err := log.NewFilterByLevel(logger, cfg.LogLevel)
	if err != nil {
		return err
	}

	// Create & start tepleton node
	privValidator := types.LoadOrGenPrivValidator(cfg.PrivValidatorFile(), tmLogger)
	n := node.NewNode(cfg, privValidator, proxy.NewLocalClientCreator(basecoinApp), tmLogger.With("module", "node"))

	_, err = n.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
