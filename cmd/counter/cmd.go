package main

import (
	"github.com/tepleton/basecoin/cmd/commands"
	"github.com/tepleton/basecoin/plugins/counter"
	"github.com/tepleton/basecoin/types"
)

func init() {
	commands.RegisterStartPlugin("counter", func() types.Plugin { return counter.New() })
}
