package cool

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth"
	bank "github.com/tepleton/tepleton-sdk/x/bank"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey
}

func TestCoolKeeper(t *testing.T) {
	ms, capKey := setupMultiStore()
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)

	am := auth.NewAccountMapper(cdc, capKey, &auth.BaseAccount{})
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, nil)
	ck := bank.NewKeeper(am)
	keeper := NewKeeper(capKey, ck, DefaultCodespace)

	err := InitGenesis(ctx, keeper, Genesis{"icy"})
	assert.Nil(t, err)

	genesis := WriteGenesis(ctx, keeper)
	assert.Nil(t, err)
	assert.Equal(t, genesis, Genesis{"icy"})

	res := keeper.GetTrend(ctx)
	assert.Equal(t, res, "icy")

	keeper.setTrend(ctx, "fiery")
	res = keeper.GetTrend(ctx)
	assert.Equal(t, res, "fiery")
}
