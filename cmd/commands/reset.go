package commands

import (
	"os"
	"path"

	"github.com/urfave/cli"

	tmcfg "github.com/tepleton/tepleton/config/tepleton"
	types "github.com/tepleton/tepleton/types"
)

var UnsafeResetAllCmd = cli.Command{
	Name:      "unsafe_reset_all",
	Usage:     "Reset all blockchain data",
	ArgsUsage: "",
	Action: func(c *cli.Context) error {
		return cmdUnsafeResetAll(c)
	},
}

func cmdUnsafeResetAll(c *cli.Context) error {
	basecoinDir := BasecoinRoot("")
	tmDir := path.Join(basecoinDir, "tepleton")
	tmConfig := tmcfg.GetConfig(tmDir)

	// Get and Reset PrivValidator
	var privValidator *types.PrivValidator
	privValidatorFile := tmConfig.GetString("priv_validator_file")
	if _, err := os.Stat(privValidatorFile); err == nil {
		privValidator = types.LoadPrivValidator(privValidatorFile)
		privValidator.Reset()
		log.Notice("Reset PrivValidator", "file", privValidatorFile)
	} else {
		privValidator = types.GenPrivValidator()
		privValidator.SetFile(privValidatorFile)
		privValidator.Save()
		log.Notice("Generated PrivValidator", "file", privValidatorFile)
	}

	// Remove all tepleton data
	tmDataDir := tmConfig.GetString("db_dir")
	os.RemoveAll(tmDataDir)
	log.Notice("Removed Tendermint data", "dir", tmDataDir)

	// Remove all basecoin data
	basecoinDataDir := path.Join(basecoinDir, "merkleeyes.db")
	os.RemoveAll(basecoinDataDir)
	log.Notice("Removed Basecoin data", "dir", basecoinDataDir)

	return nil
}
