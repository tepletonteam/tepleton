package pow

import (
	"testing"

	"github.com/stretchr/testify/require"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"

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
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	config := NewConfig("pow", int64(1))
	ck := bank.NewKeeper(am)
	keeper := NewKeeper(capKey, config, ck, DefaultCodespace)

	err := InitGenesis(ctx, keeper, Genesis{uint64(1), uint64(0)})
	require.Nil(t, err)

	genesis := WriteGenesis(ctx, keeper)
	require.Nil(t, err)
	require.Equal(t, genesis, Genesis{uint64(1), uint64(0)})

	res, err := keeper.GetLastDifficulty(ctx)
	require.Nil(t, err)
	require.Equal(t, res, uint64(1))

	keeper.SetLastDifficulty(ctx, 2)

	res, err = keeper.GetLastDifficulty(ctx)
	require.Nil(t, err)
	require.Equal(t, res, uint64(2))
}
