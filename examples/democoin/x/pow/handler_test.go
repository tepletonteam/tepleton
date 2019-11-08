package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
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
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, nil, log.NewNopLogger(), 1000)
	config := NewConfig("pow", int64(1))
	ck := bank.NewKeeper(am)
	keeper := NewKeeper(capKey, config, ck, DefaultCodespace)

	handler := keeper.Handler

	addr := sdk.Address([]byte("sender"))
	count := uint64(1)
	difficulty := uint64(2)

	err := InitGenesis(ctx, keeper, Genesis{uint64(1), uint64(0)})
	assert.Nil(t, err)

	nonce, proof := mine(addr, count, difficulty)
	msg := NewMsgMine(addr, difficulty, count, nonce, proof)

	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	newDiff, err := keeper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, newDiff, uint64(2))

	newCount, err := keeper.GetLastCount(ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCount, uint64(1))

	// todo assert correct coin change, awaiting https://github.com/tepleton/tepleton-sdk/pull/691

	difficulty = uint64(4)
	nonce, proof = mine(addr, count, difficulty)
	msg = NewMsgMine(addr, difficulty, count, nonce, proof)

	result = handler(ctx, msg)
	assert.NotEqual(t, result, sdk.Result{})
}
