package assoc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/examples/democoin/mock"
	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, wrsp.Header{}, false, nil)
	return ctx
}

func TestValidatorSet(t *testing.T) {
	key := sdk.NewKVStoreKey("test")
	ctx := defaultContext(key)

	addr1 := []byte("addr1")
	addr2 := []byte("addr2")

	base := &mock.ValidatorSet{[]mock.Validator{
		{addr1, sdk.NewRat(1)},
		{addr2, sdk.NewRat(2)},
	}}

	valset := NewValidatorSet(wire.NewCodec(), sdk.NewPrefixStoreGetter(key, []byte("assoc")), base, 1, 5)

	assert.Equal(t, base.Validator(ctx, addr1), valset.Validator(ctx, addr1))
	assert.Equal(t, base.Validator(ctx, addr2), valset.Validator(ctx, addr2))

	assoc1 := []byte("asso1")
	assoc2 := []byte("asso2")

	assert.True(t, valset.Associate(ctx, addr1, assoc1))
	assert.True(t, valset.Associate(ctx, addr2, assoc2))

	assert.Equal(t, base.Validator(ctx, addr1), valset.Validator(ctx, assoc1))
	assert.Equal(t, base.Validator(ctx, addr2), valset.Validator(ctx, assoc2))

	assert.Equal(t, base.Validator(ctx, addr1), valset.Validator(ctx, addr1))
	assert.Equal(t, base.Validator(ctx, addr2), valset.Validator(ctx, addr2))

	assocs := valset.Associations(ctx, addr1)
	assert.Equal(t, 1, len(assocs))
	assert.True(t, bytes.Equal(assoc1, assocs[0]))

	assert.False(t, valset.Associate(ctx, addr1, assoc2))
	assert.False(t, valset.Associate(ctx, addr2, assoc1))

	valset.Dissociate(ctx, addr1, assoc1)
	valset.Dissociate(ctx, addr2, assoc2)

	assert.Equal(t, base.Validator(ctx, addr1), valset.Validator(ctx, addr1))
	assert.Equal(t, base.Validator(ctx, addr2), valset.Validator(ctx, addr2))

	assert.Nil(t, valset.Validator(ctx, assoc1))
	assert.Nil(t, valset.Validator(ctx, assoc2))
}
