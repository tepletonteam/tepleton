package keys

import (
	amino "github.com/tepleton/go-amino"
	crypto "github.com/tepleton/go-crypto"
)

var cdc = amino.NewCodec()

func init() {
	crypto.RegisterAmino(cdc)
}
