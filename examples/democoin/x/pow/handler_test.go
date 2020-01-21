package pow

import (
	"testing"

	"github.com/stretchr/testify/require"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/log"

	sdk "github.com/tepleton/tepleton-sdk/types"
	wire "github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth"
	bank "github.com/tepleton/tepleton-sdk/x/bank"
)

func TestPowHandler(t *testing.T) {
	ms, capKey := setupMultiStore()
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)

	am := auth.NewAccountMapper(cdc, capKey, &auth.BaseAccount{})
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	config := NewConfig("pow", int64(1))
	ck := bank.NewKeeper(am)
	keeper := NewKeeper(capKey, config, ck, DefaultCodespace)

	handler := keeper.Handler

	addr := sdk.Address([]byte("sender"))
	count := uint64(1)
	difficulty := uint64(2)

	err := InitGenesis(ctx, keeper, Genesis{uint64(1), uint64(0)})
	require.Nil(t, err)

	nonce, proof := mine(addr, count, difficulty)
	msg := NewMsgMine(addr, difficulty, count, nonce, proof)

	result := handler(ctx, msg)
	require.Equal(t, result, sdk.Result{})

	newDiff, err := keeper.GetLastDifficulty(ctx)
	require.Nil(t, err)
	require.Equal(t, newDiff, uint64(2))

	newCount, err := keeper.GetLastCount(ctx)
	require.Nil(t, err)
	require.Equal(t, newCount, uint64(1))

	// todo assert correct coin change, awaiting https://github.com/tepleton/tepleton-sdk/pull/691

	difficulty = uint64(4)
	nonce, proof = mine(addr, count, difficulty)
	msg = NewMsgMine(addr, difficulty, count, nonce, proof)

	result = handler(ctx, msg)
	require.NotEqual(t, result, sdk.Result{})
}
