package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton/crypto"

	sdk "github.com/tepleton/tepleton-sdk/types"
)

// func newStdFee() StdFee {
// 	return NewStdFee(100,
// 		Coin{"atom", 150},
// 	)
// }

func TestStdTx(t *testing.T) {
	priv := crypto.GenPrivKeyEd25519()
	addr := priv.PubKey().Address()
	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
	fee := newStdFee()
	sigs := []StdSignature{}

	tx := NewStdTx(msgs, fee, sigs, "")
	require.Equal(t, msgs, tx.GetMsgs())
	require.Equal(t, sigs, tx.GetSignatures())

	feePayer := FeePayer(tx)
	require.Equal(t, addr, feePayer)
}
