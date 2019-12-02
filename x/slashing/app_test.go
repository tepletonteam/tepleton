package slashing

import (
	"testing"

	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/auth/mock"
	"github.com/tepleton/tepleton-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton-sdk/x/stake"
	wrsp "github.com/tepleton/wrsp/types"
	crypto "github.com/tepleton/go-crypto"
)

var (
	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	coins = sdk.Coins{sdk.NewCoin("foocoin", 10)}
)

// initialize the mock application for this module
func getMockApp(t *testing.T) (*mock.App, stake.Keeper, Keeper) {
	mapp := mock.NewApp()

	RegisterWire(mapp.Cdc)
	keyStake := sdk.NewKVStoreKey("stake")
	keySlashing := sdk.NewKVStoreKey("slashing")
	coinKeeper := bank.NewKeeper(mapp.AccountMapper)
	stakeKeeper := stake.NewKeeper(mapp.Cdc, keyStake, coinKeeper, mapp.RegisterCodespace(stake.DefaultCodespace))
	keeper := NewKeeper(mapp.Cdc, keySlashing, stakeKeeper, mapp.RegisterCodespace(DefaultCodespace))
	mapp.Router().AddRoute("stake", stake.NewHandler(stakeKeeper))
	mapp.Router().AddRoute("slashing", NewHandler(keeper))

	mapp.SetEndBlocker(getEndBlocker(stakeKeeper))
	mapp.SetInitChainer(getInitChainer(mapp, stakeKeeper))
	require.NoError(t, mapp.CompleteSetup([]*sdk.KVStoreKey{keyStake, keySlashing}))

	return mapp, stakeKeeper, keeper
}

// stake endblocker
func getEndBlocker(keeper stake.Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req wrsp.RequestEndBlock) wrsp.ResponseEndBlock {
		validatorUpdates := stake.EndBlocker(ctx, keeper)
		return wrsp.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
}

// overwrite the mock init chainer
func getInitChainer(mapp *mock.App, keeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req wrsp.RequestInitChain) wrsp.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		stake.InitGenesis(ctx, keeper, stake.DefaultGenesisState())
		return wrsp.ResponseInitChain{}
	}
}

func checkValidator(t *testing.T, mapp *mock.App, keeper stake.Keeper,
	addr sdk.Address, expFound bool) stake.Validator {
	ctxCheck := mapp.BaseApp.NewContext(true, wrsp.Header{})
	validator, found := keeper.GetValidator(ctxCheck, addr1)
	assert.Equal(t, expFound, found)
	return validator
}

func checkValidatorSigningInfo(t *testing.T, mapp *mock.App, keeper Keeper,
	addr sdk.Address, expFound bool) ValidatorSigningInfo {
	ctxCheck := mapp.BaseApp.NewContext(true, wrsp.Header{})
	signingInfo, found := keeper.getValidatorSigningInfo(ctxCheck, addr)
	assert.Equal(t, expFound, found)
	return signingInfo
}

func TestSlashingMsgs(t *testing.T) {
	mapp, stakeKeeper, keeper := getMockApp(t)

	genCoin := sdk.NewCoin("steak", 42)
	bondCoin := sdk.NewCoin("steak", 10)

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	accs := []auth.Account{acc1}
	mock.SetGenesis(mapp, accs)
	description := stake.NewDescription("foo_moniker", "", "", "")
	createValidatorMsg := stake.NewMsgCreateValidator(
		addr1, priv1.PubKey(), bondCoin, description,
	)
	mock.SignCheckDeliver(t, mapp.BaseApp, createValidatorMsg, []int64{0}, []int64{0}, true, priv1)
	mock.CheckBalance(t, mapp, addr1, sdk.Coins{genCoin.Minus(bondCoin)})
	mapp.BeginBlock(wrsp.RequestBeginBlock{})

	validator := checkValidator(t, mapp, stakeKeeper, addr1, true)
	require.Equal(t, addr1, validator.Owner)
	require.Equal(t, sdk.Bonded, validator.Status())
	require.True(sdk.RatEq(t, sdk.NewRat(10), validator.PoolShares.Bonded()))
	unrevokeMsg := MsgUnrevoke{ValidatorAddr: validator.PubKey.Address()}

	// no signing info yet
	checkValidatorSigningInfo(t, mapp, keeper, addr1, false)

	// unrevoke should fail with unknown validator
	res := mock.SignCheck(t, mapp.BaseApp, unrevokeMsg, []int64{0}, []int64{1}, priv1)
	require.Equal(t, sdk.ToWRSPCode(DefaultCodespace, CodeInvalidValidator), res.Code)
}
