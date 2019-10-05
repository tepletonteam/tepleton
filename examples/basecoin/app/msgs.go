package app

import (
	"github.com/tepleton/tepleton-sdk/x/bank"
	crypto "github.com/tepleton/go-crypto"
	wire "github.com/tepleton/go-wire"
)

// Wire requires registration of interfaces & concrete types. All
// interfaces to be encoded/decoded in a Msg must be registered
// here, along with all the concrete types that implement them.
func makeTxCodec() (cdc *wire.Codec) {
	cdc = wire.NewCodec()

	// Register crypto.[PubKey,PrivKey,Signature] types.
	crypto.RegisterWire(cdc)

	// Register bank.[SendMsg,IssueMsg] types.
	bank.RegisterWire(cdc)

	return
}
