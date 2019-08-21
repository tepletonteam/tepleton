package app

import (
	"strings"

	"github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/types"
	. "github.com/tepleton/go-common"
	"github.com/tepleton/go-wire"
	"github.com/tepleton/governmint/gov"
	eyes "github.com/tepleton/merkleeyes/client"
	wrsp "github.com/tepleton/wrsp/types"
)

const (
	version   = "0.1"
	maxTxSize = 10240

	typeByteBase = 0x01
	typeByteEyes = 0x02
	typeByteGov  = 0x03

	pluginNameBase = "base"
	pluginNameEyes = "eyes"
	pluginNameGov  = "gov"
)

type Basecoin struct {
	eyesCli *eyes.Client
	govMint *gov.Governmint
	state   *state.State
	plugins *types.Plugins
}

func NewBasecoin(eyesCli *eyes.Client) *Basecoin {
	govMint := gov.NewGovernmint(eyesCli)
	state_ := state.NewState(eyesCli)
	plugins := types.NewPlugins()
	plugins.RegisterPlugin(typeByteGov, pluginNameGov, govMint) // TODO: make constants
	return &Basecoin{
		eyesCli: eyesCli,
		govMint: govMint,
		state:   state_,
		plugins: plugins,
	}
}

// TMSP::Info
func (app *Basecoin) Info() string {
	return Fmt("Basecoin v%v", version)
}

// TMSP::SetOption
func (app *Basecoin) SetOption(key string, value string) (log string) {
	pluginName, key := splitKey(key)
	if pluginName != pluginNameBase {
		// Set option on plugin
		plugin := app.plugins.GetByName(pluginName)
		if plugin == nil {
			return "Invalid plugin name: " + pluginName
		}
		return plugin.SetOption(key, value)
	} else {
		// Set option on basecoin
		switch key {
		case "chainID":
			app.state.SetChainID(value)
			return "Success"
		case "account":
			var err error
			var acc *types.Account
			wire.ReadJSONPtr(&acc, []byte(value), &err)
			if err != nil {
				return "Error decoding acc message: " + err.Error()
			}
			app.state.SetAccount(acc.PubKey.Address(), acc)
			return "Success"
		}
		return "Unrecognized option key " + key
	}
}

// TMSP::AppendTx
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
	res = state.ExecTx(app.state, app.plugins, tx, false, nil)
	if res.IsErr() {
		return res.PrependLog("Error in AppendTx")
	}
	return wrsp.OK
}

// TMSP::CheckTx
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
	res = state.ExecTx(app.state, app.plugins, tx, true, nil)
	if res.IsErr() {
		return res.PrependLog("Error in CheckTx")
	}
	return wrsp.OK
}

// TMSP::Query
func (app *Basecoin) Query(query []byte) (res wrsp.Result) {
	if len(query) == 0 {
		return wrsp.ErrEncodingError.SetLog("Query cannot be zero length")
	}
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case typeByteBase:
		return wrsp.OK.SetLog("This type of query not yet supported")
	case typeByteEyes:
		return app.eyesCli.QuerySync(query)
	case typeByteGov:
		return app.govMint.Query(query)
	}
	return wrsp.ErrBaseUnknownPlugin.SetLog(
		Fmt("Unknown plugin with type byte %X", typeByte))
}

// TMSP::Commit
func (app *Basecoin) Commit() (res wrsp.Result) {
	// First, commit all the plugins
	for _, plugin := range app.plugins.GetList() {
		res = plugin.Commit()
		if res.IsErr() {
			PanicSanity(Fmt("Error committing plugin %v", plugin.Name))
		}
	}
	// Then, commit eyes.
	res = app.eyesCli.CommitSync()
	if res.IsErr() {
		PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// TMSP::InitChain
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	app.govMint.InitChain(validators)
}

// TMSP::EndBlock
func (app *Basecoin) EndBlock(height uint64) []*wrsp.Validator {
	app.state.ResetCacheState()
	return app.govMint.EndBlock(height)
}

//----------------------------------------

// Splits the string at the first :.
// if there are none, the second string is nil.
func splitKey(key string) (prefix string, sufix string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return key, ""
}
