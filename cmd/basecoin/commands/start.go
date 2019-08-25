package commands

import (
	"errors"
	"os"
	"path"

	"github.com/urfave/cli"

	"github.com/tepleton/wrsp/server"
	cmn "github.com/tepleton/go-common"
	cfg "github.com/tepleton/go-config"
	//logger "github.com/tepleton/go-logger"
	eyes "github.com/tepleton/merkleeyes/client"

	tmcfg "github.com/tepleton/tepleton/config/tepleton"
	"github.com/tepleton/tepleton/node"
	"github.com/tepleton/tepleton/proxy"
	tmtypes "github.com/tepleton/tepleton/types"

	"github.com/tepleton/basecoin/app"
	"github.com/tepleton/basecoin/plugins/ibc"
	"github.com/tepleton/basecoin/types"
)

var config cfg.Config

const EyesCacheSize = 10000

var StartCmd = cli.Command{
	Name:      "start",
	Usage:     "Start basecoin",
	ArgsUsage: "",
	Action: func(c *cli.Context) error {
		return cmdStart(c)
	},
	Flags: []cli.Flag{
		AddrFlag,
		EyesFlag,
		DirFlag,
		InProcTMFlag,
		ChainIDFlag,
		IbcPluginFlag,
	},
}

type plugin struct {
	name string
	init func() types.Plugin
}

var plugins = []plugin{}

// RegisterStartPlugin is used to add another
func RegisterStartPlugin(flag cli.BoolFlag, init func() types.Plugin) {
	StartCmd.Flags = append(StartCmd.Flags, flag)
	plugins = append(plugins, plugin{name: flag.GetName(), init: init})
}

func cmdStart(c *cli.Context) error {

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if c.String("eyes") == "local" {
		eyesCli = eyes.NewLocalClient(path.Join(c.String("dir"), "merkleeyes.db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(c.String("eyes"))
		if err != nil {
			return errors.New("connect to MerkleEyes: " + err.Error())
		}
	}

	// Create Basecoin app
	basecoinApp := app.NewBasecoin(eyesCli)
	if c.Bool("ibc-plugin") {
		basecoinApp.RegisterPlugin(ibc.New())
	}

	// loop through all registered plugins and enable if desired
	for _, p := range plugins {
		if c.Bool(p.name) {
			basecoinApp.RegisterPlugin(p.init())
		}
	}

	// If genesis file exists, set key-value options
	genesisFile := path.Join(c.String("dir"), "genesis.json")
	if _, err := os.Stat(genesisFile); err == nil {
		err := basecoinApp.LoadGenesis(genesisFile)
		if err != nil {
			return errors.New(cmn.Fmt("%+v", err))
		}
	}

	if c.Bool("in-proc") {
		startTendermint(c, basecoinApp)
	} else {
		startBasecoinWRSP(c, basecoinApp)
	}

	return nil
}

func startBasecoinWRSP(c *cli.Context, basecoinApp *app.Basecoin) error {
	// Start the WRSP listener
	svr, err := server.NewServer(c.String("address"), "socket", basecoinApp)
	if err != nil {
		return errors.New("create listener: " + err.Error())
	}
	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil

}

func startTendermint(c *cli.Context, basecoinApp *app.Basecoin) {
	// Get configuration
	config = tmcfg.GetConfig("")
	// logger.SetLogLevel("notice") //config.GetString("log_level"))

	// parseFlags(config, args[1:]) // Command line overrides

	// Create & start tepleton node
	privValidatorFile := config.GetString("priv_validator_file")
	privValidator := tmtypes.LoadOrGenPrivValidator(privValidatorFile)
	n := node.NewNode(config, privValidator, proxy.NewLocalClientCreator(basecoinApp))

	n.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		n.Stop()
	})
}
