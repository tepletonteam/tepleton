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
func (app *Basecoin) AppendTx(txBytes []byte) (code wrsp.CodeType, result []byte, log string) {
	if len(txBytes) > maxTxSize {
		return wrsp.CodeType_BaseEncodingError, nil, "Tx size exceeds maximum"
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return wrsp.CodeType_BaseEncodingError, nil, "Error decoding tx: " + err.Error()
	}
	// Validate and exec tx
	res = state.ExecTx(app.state, app.plugins, tx, false, nil)
	if res.IsErr() {
		return res.PrependLog("Error in AppendTx")
	}
	// Store accounts
	storeAccounts(app.eyesCli, accs)
	return wrsp.CodeType_OK, nil, "Success"
}

// TMSP::CheckTx
func (app *Basecoin) CheckTx(txBytes []byte) (code wrsp.CodeType, result []byte, log string) {
	if len(txBytes) > maxTxSize {
		return wrsp.CodeType_BaseEncodingError, nil, "Tx size exceeds maximum"
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return wrsp.CodeType_BaseEncodingError, nil, "Error decoding tx: " + err.Error()
	}
	// Validate tx
	res = state.ExecTx(app.state, app.plugins, tx, true, nil)
	if res.IsErr() {
		return res.PrependLog("Error in CheckTx")
	}
	return wrsp.CodeType_OK, nil, "Success"
}

// TMSP::Query
func (app *Basecoin) Query(query []byte) (code wrsp.CodeType, result []byte, log string) {
	return wrsp.CodeType_OK, nil, ""
	value, err := app.eyesCli.GetSync(query)
	if err != nil {
		panic("Error making query: " + err.Error())
	}
	return wrsp.CodeType_OK, value, "Success"
}

// TMSP::Commit
func (app *Basecoin) Commit() (hash []byte, log string) {
	hash, log, err := app.eyesCli.CommitSync()
	if err != nil {
		panic("Error getting hash: " + err.Error())
	}
	return hash, "Success"
}

// TMSP::InitChain
func (app *Basecoin) InitChain(validators []*wrsp.Validator) {
	app.govMint.InitChain(validators)
}

// TMSP::EndBlock
func (app *Basecoin) EndBlock(height uint64) []*wrsp.Validator {
	return app.govMint.EndBlock(height)
}

//----------------------------------------

func validateTx(tx types.Tx) (code wrsp.CodeType, errStr string) {
	inputs, outputs := tx.GetInputs(), tx.GetOutputs()
	if len(inputs) == 0 {
		return wrsp.CodeType_BaseEncodingError, "Tx.Inputs length cannot be 0"
	}
	seenPubKeys := map[string]bool{}
	signBytes := tx.SignBytes()
	for _, input := range inputs {
		code, errStr = validateInput(input, signBytes)
		if errStr != "" {
			return
		}
		keyString := input.PubKey.KeyString()
		if seenPubKeys[keyString] {
			return wrsp.CodeType_BaseEncodingError, "Duplicate input pubKey"
		}
		seenPubKeys[keyString] = true
	}
	for _, output := range outputs {
		code, errStr = validateOutput(output)
		if errStr != "" {
			return
		}
		keyString := output.PubKey.KeyString()
		if seenPubKeys[keyString] {
			return wrsp.CodeType_BaseEncodingError, "Duplicate output pubKey"
		}
		seenPubKeys[keyString] = true
	}
	sumInputs, overflow := sumAmounts(inputs, nil, 0)
	if overflow {
		return wrsp.CodeType_BaseEncodingError, "Input amount overflow"
	}
	sumOutputsPlus, overflow := sumAmounts(nil, outputs, len(inputs)+len(outputs))
	if overflow {
		return wrsp.CodeType_BaseEncodingError, "Output amount overflow"
	}
	if sumInputs < sumOutputsPlus {
		return wrsp.CodeType_BaseInsufficientFees, "Insufficient fees"
	}
	return wrsp.CodeType_OK, ""
}

func validateInput(input types.Input, signBytes []byte) (code wrsp.CodeType, errStr string) {
	if input.Amount == 0 {
		return wrsp.CodeType_BaseEncodingError, "Input amount cannot be zero"
	}
	if input.PubKey == nil {
		return wrsp.CodeType_BaseEncodingError, "Input pubKey cannot be nil"
	}
	if !input.PubKey.VerifyBytes(signBytes, input.Signature) {
		return wrsp.CodeType_BaseUnauthorized, "Invalid signature"
	}
	return wrsp.CodeType_OK, ""
}

func validateOutput(output types.Output) (code wrsp.CodeType, errStr string) {
	if output.Amount == 0 {
		return wrsp.CodeType_BaseEncodingError, "Output amount cannot be zero"
	}
	if output.PubKey == nil {
		return wrsp.CodeType_BaseEncodingError, "Output pubKey cannot be nil"
	}
	return wrsp.CodeType_OK, ""
}

func sumAmounts(inputs []types.Input, outputs []types.Output, more int) (total uint64, overflow bool) {
	total = uint64(more)
	for _, input := range inputs {
		total2 := total + input.Amount
		if total2 < total {
			return 0, true
		}
		total = total2
	}
	for _, output := range outputs {
		total2 := total + output.Amount
		if total2 < total {
			return 0, true
		}
		total = total2
	}
	return total, false
}

// Returns accounts in order of types.Tx inputs and outputs
// appendTx: true if this is for AppendTx.
// TODO: create more intelligent sequence-checking.  Current impl is just for a throughput demo.
func runTx(tx types.Tx, accMap map[string]types.PubAccount, appendTx bool) (accs []types.PubAccount, code wrsp.CodeType, errStr string) {
	switch tx := tx.(type) {
	case *types.SendTx:
		return runSendTx(tx, accMap, appendTx)
	case *types.GovTx:
		return runGovTx(tx, accMap, appendTx)
	}
	return nil, wrsp.CodeType_InternalError, "Unknown transaction type"
}

func processInputsOutputs(tx types.Tx, accMap map[string]types.PubAccount, appendTx bool) (accs []types.PubAccount, code wrsp.CodeType, errStr string) {
	inputs, outputs := tx.GetInputs(), tx.GetOutputs()
	accs = make([]types.PubAccount, 0, len(inputs)+len(outputs))
	// Deduct from inputs
	// TODO refactor, duplicated code.
	for _, input := range inputs {
		var acc, ok = accMap[input.PubKey.KeyString()]
		if !ok {
			return nil, wrsp.CodeType_BaseUnknownAccount, "Input account does not exist"
		}
		if appendTx {
			if acc.Sequence != input.Sequence {
				return nil, wrsp.CodeType_BaseBadNonce, "Invalid sequence"
			}
		} else {
			if acc.Sequence > input.Sequence {
				return nil, wrsp.CodeType_BaseBadNonce, "Invalid sequence (too low)"
			}
		}
		if acc.Balance < input.Amount {
			return nil, wrsp.CodeType_BaseInsufficientFunds, "Insufficient funds"
		}
	}
	// Add to outputs
	for _, output := range outputs {
		var acc, ok = accMap[output.PubKey.KeyString()]
		if !ok {
			// Create new account if it doesn't already exist.
			acc = types.PubAccount{
				PubKey: output.PubKey,
				Account: types.Account{
					Balance: output.Amount,
				},
			}
			accMap[output.PubKey.KeyString()] = acc
			accs = append(accs, acc)
		} else {
			// Good!
			if (acc.Balance + output.Amount) < acc.Balance {
				return nil, wrsp.CodeType_InternalError, "Output balance overflow in runTx"
			}
			acc.Balance += output.Amount
			accs = append(accs, acc)
		}
	}
	return accs, wrsp.CodeType_OK, ""
}

func runSendTx(tx types.Tx, accMap map[string]types.PubAccount, appendTx bool) (accs []types.PubAccount, code wrsp.CodeType, errStr string) {
	return processInputsOutputs(tx, accMap, appendTx)
}

func runGovTx(tx *types.GovTx, accMap map[string]types.PubAccount, appendTx bool) (accs []types.PubAccount, code wrsp.CodeType, errStr string) {
	accs, code, errStr = processInputsOutputs(tx, accMap, appendTx)
	// XXX run GovTx
	return
}

// TMSP::EndBlock
func (app *Basecoin) EndBlock(height uint64) []*wrsp.Validator {
	app.state.ResetCacheState()
	return app.govMint.EndBlock(height)
	// TODO other plugins?
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
