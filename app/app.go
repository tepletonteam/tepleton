package app

import (
	"encoding/json"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	. "github.com/tepleton/go-common"
	"github.com/tepleton/go-wire"
	eyes "github.com/tepleton/merkleeyes/client"

	sm "github.com/tepleton/basecoin/state"
	"github.com/tepleton/basecoin/types"
	"github.com/tepleton/basecoin/version"
)

const (
	maxTxSize      = 10240
	PluginNameBase = "base"
)

type Basecoin struct {
	eyesCli    *eyes.Client
	state      *sm.State
	cacheState *sm.State
	plugins    *types.Plugins
}

func NewBasecoin(eyesCli *eyes.Client) *Basecoin {
	state := sm.NewState(eyesCli)
	plugins := types.NewPlugins()
	return &Basecoin{
		eyesCli:    eyesCli,
		state:      state,
		cacheState: nil,
		plugins:    plugins,
	}
}

// XXX For testing, not thread safe!
func (app *Basecoin) GetState() *sm.State {
	return app.state.CacheWrap()
}

// WRSP::Info
func (app *Basecoin) Info() wrsp.ResponseInfo {
	return wrsp.ResponseInfo{Data: Fmt("Basecoin v%v", version.Version)}
}

func (app *Basecoin) RegisterPlugin(plugin types.Plugin) {
	app.plugins.RegisterPlugin(plugin)
}

// WRSP::SetOption
func (app *Basecoin) SetOption(key string, value string) string {
	pluginName, key := splitKey(key)
	if pluginName != PluginNameBase {
		// Set option on plugin
		plugin := app.plugins.GetByName(pluginName)
		if plugin == nil {
			return "Invalid plugin name: " + pluginName
		}
		log.Notice("SetOption on plugin", "plugin", pluginName, "key", key, "value", value)
		return plugin.SetOption(app.state, key, value)
	} else {
		// Set option on basecoin
		switch key {
		case "chain_id":
			app.state.SetChainID(value)
			return "Success"
		case "account":
			var acc types.Account
			err := json.Unmarshal([]byte(value), &acc)
			if err != nil {
				return "Error decoding acc message: " + err.Error()
			}
			app.state.SetAccount(acc.PubKey.Address(), &acc)
			log.Notice("SetAccount", "addr", acc.PubKey.Address(), "acc", acc)

			return "Success"
		}
		return "Unrecognized option key " + key
	}
}

// WRSP::DeliverTx
func (app *Basecoin) DeliverTx(txBytes []byte) (res wrsp.Result) {
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
	res = sm.ExecTx(app.state, app.plugins, tx, false, nil)
	if res.IsErr() {
		return res.PrependLog("Error in DeliverTx")
	}
	return res
}

// WRSP::CheckTx
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
	res = sm.ExecTx(app.cacheState, app.plugins, tx, true, nil)
	if res.IsErr() {
		return res.PrependLog("Error in CheckTx")
	}
	return wrsp.OK
}

// WRSP::Query
func (app *Basecoin) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Log = "Query cannot be zero length"
		resQuery.Code = wrsp.CodeType_EncodingError
		return
	}

	// handle special path for account info
	if reqQuery.Path == "/account" {
		reqQuery.Path = "/key"
		reqQuery.Data = sm.AccountKey(reqQuery.Data)
	}

	resQuery, err := app.eyesCli.QuerySync(reqQuery)
	if err != nil {
		resQuery.Log = "Failed to query MerkleEyes: " + err.Error()
		resQuery.Code = wrsp.CodeType_InternalError
		return
	}
	return
}

// WRSP::Commit
func (app *Basecoin) Commit() (res wrsp.Result) {

	// Commit state
	res = app.state.Commit()

	// Wrap the committed state in cache for CheckTx
	app.cacheState = app.state.CacheWrap()

	if res.IsErr() {
		PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// WRSP::InitChain
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	for _, plugin := range app.plugins.GetList() {
		plugin.InitChain(app.state, validators)
	}
}

// WRSP::BeginBlock
func (app *Basecoin) BeginBlock(hash []byte, header *wrsp.Header) {
	for _, plugin := range app.plugins.GetList() {
		plugin.BeginBlock(app.state, hash, header)
	}
}

// WRSP::EndBlock
func (app *Basecoin) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	for _, plugin := range app.plugins.GetList() {
		pluginRes := plugin.EndBlock(app.state, height)
		res.Diffs = append(res.Diffs, pluginRes.Diffs...)
	}
	return
}

//----------------------------------------

// Splits the string at the first '/'.
// if there are none, the second string is nil.
func splitKey(key string) (prefix string, suffix string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return key, ""
}
