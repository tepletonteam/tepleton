package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/libs/log"

	sdk "github.com/tepleton/tepleton-sdk/types"
	wire "github.com/tepleton/tepleton-sdk/wire"
)

var (
	emptyCoins = sdk.Coins{}
	oneCoin    = sdk.Coins{sdk.NewCoin("foocoin", 1)}
	twoCoins   = sdk.Coins{sdk.NewCoin("foocoin", 2)}
)

func TestFeeCollectionKeeperGetSet(t *testing.T) {
	ms, _, capKey2 := setupMultiStore()
	cdc := wire.NewCodec()

	// make context and keeper
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	fck := NewFeeCollectionKeeper(cdc, capKey2)

	// no coins initially
	currFees := fck.GetCollectedFees(ctx)
	require.True(t, currFees.IsEqual(emptyCoins))

	// set feeCollection to oneCoin
	fck.setCollectedFees(ctx, oneCoin)

	// check that it is equal to oneCoin
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(oneCoin))
}

func TestFeeCollectionKeeperAdd(t *testing.T) {
	ms, _, capKey2 := setupMultiStore()
	cdc := wire.NewCodec()

	// make context and keeper
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	fck := NewFeeCollectionKeeper(cdc, capKey2)

	// no coins initially
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(emptyCoins))

	// add oneCoin and check that pool is now oneCoin
	fck.addCollectedFees(ctx, oneCoin)
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(oneCoin))

	// add oneCoin again and check that pool is now twoCoins
	fck.addCollectedFees(ctx, oneCoin)
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(twoCoins))
}

func TestFeeCollectionKeeperClear(t *testing.T) {
	ms, _, capKey2 := setupMultiStore()
	cdc := wire.NewCodec()

	// make context and keeper
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	fck := NewFeeCollectionKeeper(cdc, capKey2)

	// set coins initially
	fck.setCollectedFees(ctx, twoCoins)
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(twoCoins))

	// clear fees and see that pool is now empty
	fck.ClearCollectedFees(ctx)
	require.True(t, fck.GetCollectedFees(ctx).IsEqual(emptyCoins))
}
