package app

import (
	"os"
	"testing"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/types"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/stretchr/testify/require"
	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/crypto"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"
)

func setGenesis(baseApp *BasecoinApp, accounts ...*types.AppAccount) (types.GenesisState, error) {
	genAccts := make([]*types.GenesisAccount, len(accounts))
	for i, appAct := range accounts {
		genAccts[i] = types.NewGenesisAccount(appAct)
	}

	genesisState := types.GenesisState{Accounts: genAccts}
	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genesisState)
	if err != nil {
		return types.GenesisState{}, err
	}

	// initialize and commit the chain
	baseApp.InitChain(wrsp.RequestInitChain{
		Validators: []wrsp.Validator{}, AppStateBytes: stateBytes,
	})
	baseApp.Commit()

	return genesisState, nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	baseApp := NewBasecoinApp(logger, db)

	// construct a pubkey and an address for the test account
	pubkey := crypto.GenPrivKeyEd25519().PubKey()
	addr := pubkey.Address()

	// construct some test coins
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)

	// create an auth.BaseAccount for the given test account and set it's coins
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	err = baseAcct.SetCoins(coins)
	require.Nil(t, err)

	// create a new test AppAccount with the given auth.BaseAccount
	appAcct := types.NewAppAccount("foobar", baseAcct)
	genState, err := setGenesis(baseApp, appAcct)
	require.Nil(t, err)

	// create a context for the BaseApp
	ctx := baseApp.BaseApp.NewContext(true, wrsp.Header{})
	res := baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)

	// reload app and ensure the account is still there
	baseApp = NewBasecoinApp(logger, db)

	stateBytes, err := wire.MarshalJSONIndent(baseApp.cdc, genState)
	require.Nil(t, err)

	// initialize the chain with the expected genesis state
	baseApp.InitChain(wrsp.RequestInitChain{
		Validators: []wrsp.Validator{}, AppStateBytes: stateBytes,
	})

	ctx = baseApp.BaseApp.NewContext(true, wrsp.Header{})
	res = baseApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, appAcct, res)
}
