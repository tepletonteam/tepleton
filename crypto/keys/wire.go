package keys

import (
	"github.com/tepleton/go-crypto"
	"github.com/tepleton/go-wire"
)

var cdc = wire.NewCodec()

func init() {
	crypto.RegisterWire(cdc)
}
