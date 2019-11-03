package client

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	bank "github.com/tepleton/tepleton-sdk/x/bank"
)

// build the sendTx msg
func BuildMsg(from sdk.Address, to sdk.Address, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}
