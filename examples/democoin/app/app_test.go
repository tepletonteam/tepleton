package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton-sdk/examples/democoin/types"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"

	wrsp "github.com/tepleton/wrsp/types"
	crypto "github.com/tepleton/go-crypto"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"
)

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

	genesisState := map[string]interface{}{
		"accounts": []*types.GenesisAccount{
			types.NewGenesisAccount(acc),
		},
		"cool": map[string]string{
			"trend": "ice-cold",
		},
	}
	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")

	vals := []wrsp.Validator{}
	bapp.InitChain(wrsp.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	bapp.Commit()

	// A checkTx context
	ctx := bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 := bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, acc, res1)

	// reload app and ensure the account is still there
	bapp = NewDemocoinApp(logger, db)
	bapp.InitChain(wrsp.RequestInitChain{AppStateBytes: []byte("{}")})
	ctx = bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 = bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, acc, res1)
}
