package mock

import (
	"testing"

	"github.com/tepleton/tepleton-sdk/baseapp"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/stretchr/testify/require"
	wrsp "github.com/tepleton/tepleton/wrsp/types"
	"github.com/tepleton/tepleton/crypto"
)

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *App, addr sdk.Address, exp sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, wrsp.Header{})
	res := app.AccountMapper.GetAccount(ctxCheck, addr)

	require.Equal(t, exp, res.GetCoins())
}

// CheckGenTx checks a generated signed transaction. The result of the check is
// compared against the parameter 'expPass'. A test assertion is made using the
// parameter 'expPass' against the result. A corresponding result is returned.
func CheckGenTx(
	t *testing.T, app *baseapp.BaseApp, msgs []sdk.Msg, accNums []int64,
	seq []int64, expPass bool, priv ...crypto.PrivKey,
) sdk.Result {
	tx := GenTx(msgs, accNums, seq, priv...)
	res := app.Check(tx)

	if expPass {
		require.Equal(t, sdk.WRSPCodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.WRSPCodeOK, res.Code, res.Log)
	}

	return res
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, app *baseapp.BaseApp, msgs []sdk.Msg, accNums []int64,
	seq []int64, expPass bool, priv ...crypto.PrivKey,
) sdk.Result {
	tx := GenTx(msgs, accNums, seq, priv...)
	res := app.Check(tx)

	if expPass {
		require.Equal(t, sdk.WRSPCodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.WRSPCodeOK, res.Code, res.Log)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(wrsp.RequestBeginBlock{})
	res = app.Deliver(tx)

	if expPass {
		require.Equal(t, sdk.WRSPCodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.WRSPCodeOK, res.Code, res.Log)
	}

	app.EndBlock(wrsp.RequestEndBlock{})
	app.Commit()

	return res
}
