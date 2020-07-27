package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton-sdk/examples/democoin/types"
	"github.com/tepleton/tepleton-sdk/examples/democoin/x/cool"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/crypto"
	dbm "github.com/tepleton/tepleton/libs/db"
	"github.com/tepleton/tepleton/libs/log"
)

func setGenesis(bapp *DemocoinApp, trend string, accs ...auth.BaseAccount) error {
	genaccs := make([]*types.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = types.NewGenesisAccount(&types.AppAccount{acc, "foobart"})
	}

	genesisState := types.GenesisState{
		Accounts:    genaccs,
		CoolGenesis: cool.Genesis{trend},
	}

	stateBytes, err := wire.MarshalJSONIndent(bapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []wrsp.Validator{}
	bapp.InitChain(wrsp.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	bapp.Commit()

	return nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	bapp := NewDemocoinApp(logger, db)

	// Construct some genesis bytes to reflect democoin/types/AppAccount
	pk := crypto.GenPrivKeyEd25519().PubKey()
	addr := pk.Address()
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)
	baseAcc := auth.BaseAccount{
		Address: addr,
		Coins:   coins,
	}
	acc := &types.AppAccount{baseAcc, "foobart"}

	err = setGenesis(bapp, "ice-cold", baseAcc)
	require.Nil(t, err)
	// A checkTx context
	ctx := bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 := bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	require.Equal(t, acc, res1)

	// reload app and ensure the account is still there
	bapp = NewDemocoinApp(logger, db)
	bapp.InitChain(wrsp.RequestInitChain{AppStateBytes: []byte("{}")})
	ctx = bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 = bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	require.Equal(t, acc, res1)
}
