package wire

import (
	"bytes"
	"reflect"

	"github.com/tepleton/go-amino"
)

type Codec = amino.Codec

func NewCodec() *Codec {
	cdc := amino.NewCodec()
	RegisterAmino(cdc)
	return cdc
}
