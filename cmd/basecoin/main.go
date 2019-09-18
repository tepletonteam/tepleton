package main

import (
	"os"

	"github.com/tepleton/tmlibs/cli"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/cmd/basecoin/commands"
	"github.com/tepleton/basecoin/modules/auth"
	"github.com/tepleton/basecoin/modules/base"
	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/modules/fee"
	"github.com/tepleton/basecoin/modules/abi"
	"github.com/tepleton/basecoin/modules/nonce"
	"github.com/tepleton/basecoin/modules/roles"
	"github.com/tepleton/basecoin/stack"
)

// BuildApp constructs the stack we want to use for this app
func BuildApp(feeDenom string) basecoin.Handler {
	// use the default stack
	c := coin.NewHandler()
	r := roles.NewHandler()
	i := abi.NewHandler()

	return stack.New(
		base.Logger{},
		stack.Recovery{},
		auth.Signatures{},
		base.Chain{},
		stack.Checkpoint{OnCheck: true},
		nonce.ReplayCheck{},
	).
		ABI(abi.NewMiddleware()).
		Apps(
			roles.NewMiddleware(),
			fee.NewSimpleFeeMiddleware(coin.Coin{feeDenom, 0}, fee.Bank),
			stack.Checkpoint{OnDeliver: true},
		).
		Dispatch(
			stack.WrapHandler(c),
			stack.WrapHandler(r),
			stack.WrapHandler(i),
		)
}

func main() {
	rt := commands.RootCmd

	// require all fees in mycoin - change this in your app!
	commands.Handler = BuildApp("mycoin")

	rt.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		//commands.RelayCmd,
		commands.UnsafeResetAllCmd,
		commands.VersionCmd,
	)

	cmd := cli.PrepareMainCmd(rt, "BC", os.ExpandEnv("$HOME/.basecoin"))
	cmd.Execute()
}
