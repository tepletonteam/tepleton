package app

import (
	"os"
	"fmt"
	"testing"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton-sdk/examples/basecoin/types"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/stake"
	gen "github.com/tepleton/tepleton-sdk/x/stake/types"

	wrsp "github.com/tepleton/wrsp/types"
	crypto "github.com/tepleton/go-crypto"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"
)

func setGenesis(bapp *BasecoinApp, accs ...auth.BaseAccount) error {
	genaccs := make([]*types.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = types.NewGenesisAccount(&types.AppAccount{acc, "foobart"})
	}

	genesisState := types.GenesisState{
		Accounts:  genaccs,
		StakeData: stake.DefaultGenesisState(),
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

//_______________________________________________________________________

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	bapp := NewBasecoinApp(logger, db)

	// Construct some genesis bytes to reflect basecoin/types/AppAccount
	pk := crypto.GenPrivKeyEd25519().PubKey()
	addr := pk.Address()
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)
	baseAcc := auth.BaseAccount{
		Address: addr,
		Coins:   coins,
	}
	acc := &types.AppAccount{baseAcc, "foobart"}

	err = setGenesis(bapp, baseAcc)
	require.Nil(t, err)

	// A checkTx context
	ctx := bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 := bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, acc, res1)

	// reload app and ensure the account is still there
	bapp = NewBasecoinApp(logger, db)
	// Initialize stake data with default genesis state
	stakedata := gen.DefaultGenesisState()
	genState, err := json.Marshal(stakedata)
	if err != nil {
		panic(err)
	}

	// InitChain with default stake data. Initializes deliverState and checkState context
	bapp.InitChain(wrsp.RequestInitChain{AppStateBytes: []byte(fmt.Sprintf("{\"stake\": %s}", string(genState)))})
	
	ctx = bapp.BaseApp.NewContext(true, wrsp.Header{})
	res1 = bapp.accountMapper.GetAccount(ctx, baseAcc.Address)
	assert.Equal(t, acc, res1)
}
