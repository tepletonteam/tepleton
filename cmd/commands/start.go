package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tepleton/wrsp/server"
	cmn "github.com/tepleton/tmlibs/common"
	eyes "github.com/tepleton/merkleeyes/client"

	tmcfg "github.com/tepleton/tepleton/config/tepleton"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	tmtypes "github.com/tepleton/tepleton/types"

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
	basecoinDir := BasecoinRoot("")

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if eyesFlag == "local" {
		eyesCli = eyes.NewLocalClient(path.Join(basecoinDir, "data", "merkleeyes.db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(eyesFlag)
		if err != nil {
			return errors.Errorf("Error connecting to MerkleEyes: %v\n", err)
		}
	}

	// Create Basecoin app
	basecoinApp := app.NewBasecoin(eyesCli)

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
		genesisFile := path.Join(basecoinDir, "genesis.json")
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
		log.Notice("Starting Basecoin without Tendermint", "chain_id", chainID)
		// run just the wrsp app/server
		return startBasecoinWRSP(basecoinApp)
	} else {
		log.Notice("Starting Basecoin with Tendermint", "chain_id", chainID)
		// start the app with tepleton in-process
		return startTendermint(basecoinDir, basecoinApp)
	}
}

func startBasecoinWRSP(basecoinApp *app.Basecoin) error {

	// Start the WRSP listener
	svr, err := server.NewServer(addrFlag, "socket", basecoinApp)
	if err != nil {
		return errors.Errorf("Error creating listener: %v\n", err)
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}

func startTendermint(dir string, basecoinApp *app.Basecoin) error {

	// Get configuration
	tmConfig := tmcfg.GetConfig(dir)
	// logger.SetLogLevel("notice") //config.GetString("log_level"))
	// parseFlags(config, args[1:]) // Command line overrides

	// Create & start tepleton node
	privValidatorFile := tmConfig.GetString("priv_validator_file")
	privValidator := tmtypes.LoadOrGenPrivValidator(privValidatorFile)
	n := node.NewNode(tmConfig, privValidator, proxy.NewLocalClientCreator(basecoinApp))

	_, err := n.Start()
	if err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		n.Stop()
	})
	return nil
}
