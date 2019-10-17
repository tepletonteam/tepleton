package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"

	sdk "github.com/tepleton/tepleton-sdk/types"
	auth "github.com/tepleton/tepleton-sdk/x/auth"
	bank "github.com/tepleton/tepleton-sdk/x/bank"
)

func TestPowHandler(t *testing.T) {
	ms, capKey := setupMultiStore()

	am := auth.NewAccountMapper(capKey, &auth.BaseAccount{})
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, nil)
	mapper := NewMapper(capKey)
	config := NewPowConfig("pow", int64(1))
	ck := bank.NewCoinKeeper(am)

	handler := NewHandler(ck, mapper, config)

	addr := sdk.Address([]byte("sender"))
	count := uint64(1)
	difficulty := uint64(2)
	nonce, proof := mine(addr, count, difficulty)
	msg := NewMineMsg(addr, difficulty, count, nonce, proof)

	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	newDiff, err := mapper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, newDiff, uint64(2))

	newCount, err := mapper.GetLastCount(ctx)
	assert.Nil(t, err)
	assert.Equal(t, newCount, uint64(1))

	// todo assert correct coin change, awaiting https://github.com/tepleton/tepleton-sdk/pull/691

	difficulty = uint64(4)
	nonce, proof = mine(addr, count, difficulty)
	msg = NewMineMsg(addr, difficulty, count, nonce, proof)

	result = handler(ctx, msg)
	assert.NotEqual(t, result, sdk.Result{})
}
