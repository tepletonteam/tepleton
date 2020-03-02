package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/tepleton/tepleton-sdk/types"
)

func TestMsgUnrevokeGetSignBytes(t *testing.T) {
	addr := sdk.Address("abcd")
	msg := NewMsgUnrevoke(addr)
	bytes := msg.GetSignBytes()
	require.Equal(t, string(bytes), `{"address":"tepletonvaladdr1v93xxeqamr0mv"}`)
}
