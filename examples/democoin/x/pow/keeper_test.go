package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth"
	bank "github.com/tepleton/tepleton-sdk/x/bank"
)

// possibly share this kind of setup functionality between module testsuites?
func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, capKey
}

func TestPowKeeperGetSet(t *testing.T) {
	ms, capKey := setupMultiStore()
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)

	am := auth.NewAccountMapper(cdc, capKey, &auth.BaseAccount{})
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, nil)
	config := NewPowConfig("pow", int64(1))
	ck := bank.NewCoinKeeper(am)
	keeper := NewKeeper(capKey, config, ck)

	err := keeper.InitGenesis(ctx, PowGenesis{uint64(1), uint64(0)})
	assert.Nil(t, err)

	res, err := keeper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, uint64(1))

	keeper.SetLastDifficulty(ctx, 2)

	res, err = keeper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, uint64(2))
}
