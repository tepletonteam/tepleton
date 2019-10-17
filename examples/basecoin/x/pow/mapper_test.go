package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
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

func TestPowMapperGetSet(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, wrsp.Header{}, false, nil)
	mapper := NewMapper(capKey)

	res, err := mapper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, uint64(1))

	mapper.SetLastDifficulty(ctx, 2)

	res, err = mapper.GetLastDifficulty(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, uint64(2))
}
