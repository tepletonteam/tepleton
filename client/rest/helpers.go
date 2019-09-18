package rest

import (
	keycmd "github.com/tepleton/go-crypto/cmd"
	"github.com/tepleton/go-crypto/keys"
	wire "github.com/tepleton/go-wire"

	ctypes "github.com/tepleton/tepleton/rpc/core/types"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/client/commands"
)

// PostTx is same as a tx
func PostTx(tx basecoin.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	packet := wire.BinaryBytes(tx)
	// post the bytes
	node := commands.GetNode()
	return node.BroadcastTxCommit(packet)
}

// SignTx will modify the tx in-place, adding a signature if possible
func SignTx(name, pass string, tx basecoin.Tx) error {
	if sign, ok := tx.Unwrap().(keys.Signable); ok {
		manager := keycmd.GetKeyManager()
		return manager.Sign(name, pass, sign)
	}
	return nil
}
