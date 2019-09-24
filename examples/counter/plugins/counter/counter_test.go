package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-wire"
	"github.com/tepleton/tmlibs/log"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/app"
	"github.com/tepleton/tepleton-sdk/modules/auth"
	"github.com/tepleton/tepleton-sdk/modules/base"
	"github.com/tepleton/tepleton-sdk/modules/coin"
	"github.com/tepleton/tepleton-sdk/modules/nonce"
)

func TestCounterPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Basecoin initialization
	chainID := "test_chain_id"
	logger := log.TestingLogger()
	// logger := log.NewTracingLogger(log.NewTMLogger(os.Stdout))

	h := NewHandler("gold")
	store, err := app.MockStoreApp("counter", logger)
	require.Nil(err, "%+v", err)
	bcApp := app.NewBaseApp(store, h, nil)
	_, err = bcApp.InitState("base", "chain_id", chainID)
	require.Nil(err, "%+v", err)

	// Account initialization
	bal := coin.Coins{{"", 1000}, {"gold", 1000}}
	acct := coin.NewAccountWithKey(bal)
	_, err = bcApp.InitState("coin", "account", acct.MakeOption())
	require.Nil(err, "%+v", err)

	// Deliver a CounterTx
	DeliverCounterTx := func(valid bool, counterFee coin.Coins, sequence uint32) wrsp.Result {
		tx := NewTx(valid, counterFee)
		tx = nonce.NewTx(sequence, []sdk.Actor{acct.Actor()}, tx)
		tx = base.NewChainTx(chainID, 0, tx)
		stx := auth.NewSig(tx)
		auth.Sign(stx, acct.Key)
		txBytes := wire.BinaryBytes(stx.Wrap())
		return bcApp.DeliverTx(txBytes)
	}

	// Test a basic send, no fee
	res := DeliverCounterTx(true, nil, 1)
	assert.True(res.IsOK(), res.String())

	// Test an invalid send, no fee
	res = DeliverCounterTx(false, nil, 2)
	assert.True(res.IsErr(), res.String())

	// Test an invalid sequence
	res = DeliverCounterTx(true, nil, 2)
	assert.True(res.IsErr(), res.String())

	// Test an valid send, with supported fee
	res = DeliverCounterTx(true, coin.Coins{{"gold", 100}}, 3)
	assert.True(res.IsOK(), res.String())

	// Test unsupported fee
	res = DeliverCounterTx(true, coin.Coins{{"silver", 100}}, 4)
	assert.True(res.IsErr(), res.String())
}
