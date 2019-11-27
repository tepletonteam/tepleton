package rpc

import (
	amino "github.com/tepleton/go-amino"
	ctypes "github.com/tepleton/tepleton/rpc/core/types"
)

var cdc = amino.NewCodec()

func init() {
	ctypes.RegisterAmino(cdc)
}
