package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/app"
	"github.com/tepleton/basecoin/modules/auth"
	"github.com/tepleton/basecoin/modules/base"
	"github.com/tepleton/basecoin/modules/coin"
	"github.com/tepleton/basecoin/modules/nonce"
	"github.com/tepleton/basecoin/state/merkle"
	"github.com/tepleton/go-wire"
	"github.com/tepleton/tmlibs/log"
)

func TestCounterPlugin(t *testing.T) {
	assert := assert.New(t)

	// Basecoin initialization
	chainID := "test_chain_id"
	logger := log.TestingLogger()
	// logger := log.NewTracingLogger(log.NewTMLogger(os.Stdout))

	store := merkle.NewStore("", 0, logger.With("module", "store"))
	h := NewHandler("gold")
	bcApp := app.NewBasecoin(
		h,
		store,
		logger.With("module", "app"),
	)
	bcApp.SetOption("base/chain_id", chainID)

	// Account initialization
	bal := coin.Coins{{"", 1000}, {"gold", 1000}}
	acct := coin.NewAccountWithKey(bal)
	log := bcApp.SetOption("coin/account", acct.MakeOption())
	require.Equal(t, "Success", log)

	// Deliver a CounterTx
	DeliverCounterTx := func(valid bool, counterFee coin.Coins, sequence uint32) wrsp.Result {
		tx := NewTx(valid, counterFee)
		tx = nonce.NewTx(sequence, []basecoin.Actor{acct.Actor()}, tx)
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
