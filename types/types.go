package types

import (
	"github.com/tepleton/go-crypto"
)

type Tx struct {
	Inputs  []Input
	Outputs []Output
}

type Input struct {
	PubKey    crypto.PubKey
	Amount    uint64
	Sequence  uint
	Signature crypto.Signature
}

type Output struct {
	PubKey crypto.PubKey
	Amount uint64
}

type Account struct {
	Sequence uint
	Balance  uint64
	crypto.PubKey `json:"-"`
}
