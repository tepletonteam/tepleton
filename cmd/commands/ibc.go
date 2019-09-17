package commands

import "github.com/tepleton/basecoin/plugins/abi"

// returns a new ABI plugin to be registered with Basecoin
func NewABIPlugin() *abi.ABIPlugin {
	return abi.New()
}
