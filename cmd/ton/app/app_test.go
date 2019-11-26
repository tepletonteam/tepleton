package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/stake"

	wrsp "github.com/tepleton/wrsp/types"
	crypto "github.com/tepleton/go-crypto"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"
)

func setGenesis(gapp *GaiaApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	genesisState := GenesisState{
		Accounts:  genaccs,
		StakeData: stake.DefaultGenesisState(),
	}

	stateBytes, err := wire.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []wrsp.Validator{}
	gapp.InitChain(wrsp.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	gapp.Commit()

	return nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	gapp := NewGaiaApp(logger, db)

	// Construct some genesis bytes to reflect GaiaAccount
	pk := crypto.GenPrivKeyEd25519().PubKey()
	addr := pk.Address()
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)
	baseAcc := &auth.BaseAccount{
		Address: addr,
		Coins:   coins,
	}

	err = setGenesis(gapp, baseAcc)
	require.Nil(t, err)

	// A checkTx context
	ctx := gapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 := gapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, baseAcc, res1)

	// reload app and ensure the account is still there
	gapp = NewGaiaApp(logger, db)
	ctx = gapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 = gapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, baseAcc, res1)
}
