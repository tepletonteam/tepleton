package wire

import (
	"github.com/tepleton/go-amino"
	"github.com/tepleton/go-crypto"
)

type Codec = amino.Codec

func NewCodec() *Codec {
	cdc := amino.NewCodec()
	return cdc
}

func RegisterCrypto(cdc *Codec) {
	crypto.RegisterAmino(cdc)
}
