package ibc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-crypto"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, wrsp.Header{}, false, nil)
	return ctx
}

func newAddress() crypto.Address {
	return crypto.GenPrivKeyEd25519().PubKey().Address()
}

func TestHandler(t *testing.T) {
	key := sdk.NewKVStoreKey("ibc")
	ctx := defaultContext(key)
	ibcm := NewIBCMapper(key)

	h := NewHandler(ibcm)

	src := newAddress()
	dest := newAddress()
	chainid := "ibcchain"
	coin := sdk.Coin{Denom: "neutron", Amount: 10000}

	packet := IBCPacket{
		SrcAddr:   src,
		DestAddr:  dest,
		Coins:     sdk.Coins{coin},
		SrcChain:  chainid,
		DestChain: chainid,
	}

	store := ctx.KVStore(key)

	var msg sdk.Msg
	var res sdk.Result
	var egl int64
	var igs int64

	egl = ibcm.getEgressLength(store, chainid)
	assert.Equal(t, egl, int64(0))

	msg = IBCTransferMsg{
		IBCPacket: packet,
	}
	res = h(ctx, msg)
	assert.True(t, res.IsOK())

	egl = ibcm.getEgressLength(store, chainid)
	assert.Equal(t, egl, int64(1))

	igs = ibcm.GetIngressSequence(ctx, chainid)
	assert.Equal(t, igs, int64(0))

	msg = IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   src,
		Sequence:  0,
	}
	res = h(ctx, msg)
	assert.True(t, res.IsOK())

	igs = ibcm.GetIngressSequence(ctx, chainid)
	assert.Equal(t, igs, int64(1))

	res = h(ctx, msg)
	assert.False(t, res.IsOK())

	igs = ibcm.GetIngressSequence(ctx, chainid)
	assert.Equal(t, igs, int64(1))
}
