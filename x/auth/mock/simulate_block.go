package mock

import (
	"testing"

	"github.com/tepleton/tepleton-sdk/baseapp"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tepleton/go-crypto"

	wrsp "github.com/tepleton/wrsp/types"
)

var chainID = "" // TODO

// set the mock app genesis
func SetGenesis(app *App, accs []auth.Account) {

	// pass the accounts in via the application (lazy) instead of through RequestInitChain
	app.GenesisAccounts = accs

	app.InitChain(wrsp.RequestInitChain{})
	app.Commit()
}

// check an account balance
func CheckBalance(t *testing.T, app *App, addr sdk.Address, exp sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, wrsp.Header{})
	res := app.AccountMapper.GetAccount(ctxCheck, addr)
	assert.Equal(t, exp, res.GetCoins())
}

// generate a signed transaction
func GenTx(msg sdk.Msg, seq []int64, priv ...crypto.PrivKeyEd25519) auth.StdTx {

	// make the transaction free
	fee := auth.StdFee{
		sdk.Coins{{"foocoin", 0}},
		100000,
	}

	sigs := make([]auth.StdSignature, len(priv))
	for i, p := range priv {
		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: p.Sign(auth.StdSignBytes(chainID, seq, fee, msg)),
			Sequence:  seq[i],
		}
	}
	return auth.NewStdTx(msg, fee, sigs)
}

// check a transaction result
func SignCheck(t *testing.T, app *baseapp.BaseApp, msg sdk.Msg, seq []int64, priv ...crypto.PrivKeyEd25519) sdk.Result {
	tx := GenTx(msg, seq, priv...)
	res := app.Check(tx)
	return res
}

// simulate a block
func SignCheckDeliver(t *testing.T, app *baseapp.BaseApp, msg sdk.Msg, seq []int64, expPass bool, priv ...crypto.PrivKeyEd25519) {

	// Sign the tx
	tx := GenTx(msg, seq, priv...)

	// Run a Check
	res := app.Check(tx)
	if expPass {
		require.Equal(t, sdk.WRSPCodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.WRSPCodeOK, res.Code, res.Log)
	}

	// Simulate a Block
	app.BeginBlock(wrsp.RequestBeginBlock{})
	res = app.Deliver(tx)
	if expPass {
		require.Equal(t, sdk.WRSPCodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.WRSPCodeOK, res.Code, res.Log)
	}
	app.EndBlock(wrsp.RequestEndBlock{})

	app.Commit()
}
