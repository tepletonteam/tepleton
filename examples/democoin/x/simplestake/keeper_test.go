package simplestake

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/crypto"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/bank"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey, capKey
}

func TestKeeperGetSet(t *testing.T) {
	ms, authKey, capKey := setupMultiStore()
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)

	accountMapper := auth.NewAccountMapper(cdc, authKey, &auth.BaseAccount{})
	stakeKeeper := NewKeeper(capKey, bank.NewKeeper(accountMapper), DefaultCodespace)
	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())
	addr := sdk.Address([]byte("some-address"))

	bi := stakeKeeper.getBondInfo(ctx, addr)
	require.Equal(t, bi, bondInfo{})

	privKey := crypto.GenPrivKeyEd25519()

	bi = bondInfo{
		PubKey: privKey.PubKey(),
		Power:  int64(10),
	}
	fmt.Printf("Pubkey: %v\n", privKey.PubKey())
	stakeKeeper.setBondInfo(ctx, addr, bi)

	savedBi := stakeKeeper.getBondInfo(ctx, addr)
	require.NotNil(t, savedBi)
	fmt.Printf("Bond Info: %v\n", savedBi)
	require.Equal(t, int64(10), savedBi.Power)
}

func TestBonding(t *testing.T) {
	ms, authKey, capKey := setupMultiStore()
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, wrsp.Header{}, false, log.NewNopLogger())

	accountMapper := auth.NewAccountMapper(cdc, authKey, &auth.BaseAccount{})
	coinKeeper := bank.NewKeeper(accountMapper)
	stakeKeeper := NewKeeper(capKey, coinKeeper, DefaultCodespace)
	addr := sdk.Address([]byte("some-address"))
	privKey := crypto.GenPrivKeyEd25519()
	pubKey := privKey.PubKey()

	_, _, err := stakeKeeper.unbondWithoutCoins(ctx, addr)
	require.Equal(t, err, ErrInvalidUnbond(DefaultCodespace))

	_, err = stakeKeeper.bondWithoutCoins(ctx, addr, pubKey, sdk.NewCoin("steak", 10))
	require.Nil(t, err)

	power, err := stakeKeeper.bondWithoutCoins(ctx, addr, pubKey, sdk.NewCoin("steak", 10))
	require.Nil(t, err)
	require.Equal(t, int64(20), power)

	pk, _, err := stakeKeeper.unbondWithoutCoins(ctx, addr)
	require.Nil(t, err)
	require.Equal(t, pubKey, pk)

	_, _, err = stakeKeeper.unbondWithoutCoins(ctx, addr)
	require.Equal(t, err, ErrInvalidUnbond(DefaultCodespace))
}
