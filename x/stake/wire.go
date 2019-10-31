package stake

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgDeclareCandidacy{}, "tepleton-sdk/MsgDeclareCandidacy", nil)
	cdc.RegisterConcrete(MsgEditCandidacy{}, "tepleton-sdk/MsgEditCandidacy", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "tepleton-sdk/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "tepleton-sdk/MsgUnbond", nil)
}
