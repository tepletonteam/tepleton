package commands

import (
	wire "github.com/tepleton/go-wire"
	"github.com/tepleton/light-client/proofs"

	"github.com/tepleton/basecoin"
)

// BaseTxPresenter this decodes all basecoin tx
type BaseTxPresenter struct {
	proofs.RawPresenter // this handles MakeKey as hex bytes
}

// ParseData - unmarshal raw bytes to a basecoin tx
func (BaseTxPresenter) ParseData(raw []byte) (interface{}, error) {
	var tx basecoin.Tx
	err := wire.ReadBinaryBytes(raw, &tx)
	return tx, err
}
