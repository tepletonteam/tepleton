package gov

import (
	"github.com/tepleton/tepleton-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {

	cdc.RegisterConcrete(MsgSubmitProposal{}, "tepleton-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "tepleton-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "tepleton-sdk/MsgVote", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)
}

var msgCdc = wire.NewCodec()
