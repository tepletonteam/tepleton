package app

import (
	"github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/types"
	. "github.com/tepleton/go-common"
	"github.com/tepleton/go-wire"
	"github.com/tepleton/governmint/gov"
	eyes "github.com/tepleton/merkleeyes/client"
	wrsp "github.com/tepleton/wrsp/types"
)

const version = "0.1"
const maxTxSize = 10240

type Basecoin struct {
	eyesCli *eyes.Client
	govMint *gov.Governmint
	state   *state.State
}

func NewBasecoin(eyesCli *eyes.Client) *Basecoin {
	govMint := gov.NewGovernmint(eyesCli)
	return &Basecoin{
		eyesCli: eyesCli,
		govMint: govMint,
		state:   state.NewState(eyesCli),
	}
}

// wrsp::Info
func (app *Basecoin) Info() string {
	return Fmt("Basecoin v%v", version)
}

// wrsp::SetOption
func (app *Basecoin) SetOption(key string, value string) (log string) {
	switch key {
	case "chainID":
		app.state.SetChainID(value)
		return "Success"
	case "account":
		var err error
		var setAccount types.Account
		wire.ReadJSONPtr(&setAccount, []byte(value), &err)
		if err != nil {
			return "Error decoding setAccount message: " + err.Error()
		}
		accBytes := wire.BinaryBytes(setAccount)
		err = app.eyesCli.SetSync(setAccount.PubKey.Address(), accBytes)
		if err != nil {
			return "Error saving account: " + err.Error()
		}
		return "Success"
	}
	return "Unrecognized option key " + key
}

// wrsp::AppendTx
func (app *Basecoin) AppendTx(txBytes []byte) (res wrsp.Result) {
	if len(txBytes) > maxTxSize {
		return wrsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return wrsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate and exec tx
	res = state.ExecTx(app.state, tx, false, nil)
	if res.IsErr() {
		return res.PrependLog("Error in AppendTx")
	}
	return wrsp.OK
}

// wrsp::CheckTx
func (app *Basecoin) CheckTx(txBytes []byte) (res wrsp.Result) {
	if len(txBytes) > maxTxSize {
		return wrsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return wrsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate tx
	res = state.ExecTx(app.state, tx, true, nil)
	if res.IsErr() {
		return res.PrependLog("Error in CheckTx")
	}
	return wrsp.OK
}

// wrsp::Query
func (app *Basecoin) Query(query []byte) (res wrsp.Result) {
	return wrsp.OK
	res = app.eyesCli.GetSync(query)
	if res.IsErr() {
		return res.PrependLog("Error querying eyesCli")
	}
	return res
}

// wrsp::Commit
func (app *Basecoin) Commit() (res wrsp.Result) {
	res = app.eyesCli.CommitSync()
	if res.IsErr() {
		panic("Error getting hash: " + res.Error())
	}
	return res
}

// wrsp::InitChain
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	app.govMint.InitChain(validators)
}

// wrsp::EndBlock
func (app *Basecoin) EndBlock(height uint64) []*wrsp.Validator {
	app.state.ResetCacheState()
	return app.govMint.EndBlock(height)
}
