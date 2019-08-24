package main

import (
	"errors"

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
	"github.com/tepleton/basecoin/plugins/counter"
	"github.com/tepleton/basecoin/plugins/ibc"
)

var config cfg.Config

const EyesCacheSize = 10000

func cmdStart(c *cli.Context) error {

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if c.String("eyes") == "local" {
		eyesCli = eyes.NewLocalClient(c.String("eyes-db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(c.String("eyes"))
		if err != nil {
			return errors.New("connect to MerkleEyes: " + err.Error())
		}
	}

	// Create Basecoin app
	basecoinApp := app.NewBasecoin(eyesCli)

	if c.Bool("counter-plugin") {
		basecoinApp.RegisterPlugin(counter.New("counter"))
	}

	if c.Bool("ibc-plugin") {
		basecoinApp.RegisterPlugin(ibc.New())

	}

	// If genesis file was specified, set key-value options
	if c.String("genesis") != "" {
		err := basecoinApp.LoadGenesis(c.String("genesis"))
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
