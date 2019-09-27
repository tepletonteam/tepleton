package app

import (
	"bytes"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	sdk "github.com/tepleton/tepleton-sdk"
)

const mainKeyHeader = "header"

// App - The WRSP application
type App struct {
	logger log.Logger

	// App name from wrsp.Info
	name string

	// DeliverTx (main) state
	store MultiStore

	// CheckTx state
	storeCheck CacheMultiStore

	// Current block header
	header *wrsp.Header

	// Handler for CheckTx and DeliverTx.
	handler sdk.Handler

	// Cached validator changes from DeliverTx
	valSetDiff []wrsp.Validator
}

var _ wrsp.Application = &App{}

func NewApp(name string) *App {
	return &App{
		name:   name,
		logger: makeDefaultLogger(),
	}
}

func (app *App) SetStore(store MultiStore) {
	app.store = store
}

func (app *App) SetHandler(handler Handler) {
	app.handler = handler
}

func (app *App) LoadLatestVersion() error {
	curVersion := app.store.NextVersion()
	app.
}

func (app *App) LoadVersion(version int64) error {
	store := app.store
	lastCommitID := store.LastCommitID()
	curVersion := store.NextVersion()
	main := store.GetKVStore("main")
	header := (*wrsp.Header)(nil)
	storeCheck := store.CacheMultiStore()

	// Main store should exist.
	if store.GetKVStore("main") == nil {
		return errors.New("App expects MultiStore with 'main' KVStore")
	}

	// Basic sanity check.
	if curVersion != lastCommitID.Version+1 {
		return errors.New("NextVersion != LastCommitID.Version+1")
	}

	// If we've committed before, we expect store(main)/<mainKeyHeader>.
	if !lastCommitID.IsZero() {
		headerBytes, ok := main.Get(mainKeyHeader)
		if !ok {
			errStr := fmt.Sprintf("Version > 0 but missing key %s", mainKeyHeader)
			return errors.New(errStr)
		}
		err = proto.Unmarshal(headerBytes, header)
		if err != nil {
			return errors.Wrap(err, "Failed to parse Header")
		}
		if header.Height != curVersion-1 {
			errStr := fmt.Sprintf("Expected header.Height %v but got %v", version, headerHeight)
			return errors.New(errStr)
		}
	}
}

//----------------------------------------

// DeliverTx - WRSP - dispatches to the handler
func (app *App) DeliverTx(txBytes []byte) wrsp.ResponseDeliverTx {

	// Initialize arguments to Handler.
	var isCheckTx = false
	var ctx = sdk.NewContext(app.header, isCheckTx, txBytes)
	var store = app.store
	var tx Tx = nil // nil until a decorator parses one.

	// Run the handler.
	var result = app.handler(ctx, app.store, tx)

	// After-handler hooks.
	// TODO move to app.afterHandler(...).
	if result.Code == wrsp.CodeType_OK {
		// XXX No longer "diff", we need to replace old entries.
		app.ValSetDiff = append(app.ValSetDiff, result.ValSetDiff)
	} else {
		// Even though the Code is not OK, there will be some side effects,
		// like those caused by fee deductions or sequence incrementations.
	}

	// Tell the blockchain engine (i.e. Tendermint).
	return wrsp.ResponseDeliverTx{
		Code: result.Code,
		Data: result.Data,
		Log:  result.Log,
		Tags: result.Tags,
	}
}

// CheckTx - WRSP - dispatches to the handler
func (app *App) CheckTx(txBytes []byte) wrsp.ResponseCheckTx {

	// Initialize arguments to Handler.
	var isCheckTx = true
	var ctx = sdk.NewContext(app.header, isCheckTx, txBytes)
	var store = app.store
	var tx Tx = nil // nil until a decorator parses one.

	// Run the handler.
	var result = app.handler(ctx, app.store, tx)

	// Tell the blockchain engine (i.e. Tendermint).
	return wrsp.ResponseDeliverTx{
		Code:      result.Code,
		Data:      result.Data,
		Log:       result.Log,
		Gas:       result.Gas,
		FeeDenom:  result.FeeDenom,
		FeeAmount: result.FeeAmount,
	}
}

// Info - WRSP
func (app *App) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.store.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.Name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// SetOption - WRSP
func (app *App) SetOption(key string, value string) string {
	return "Not Implemented"
}

// Query - WRSP
func (app *App) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	/*
		XXX Make this work with MultiStore.
		XXX It will require some interfaces updates in store/types.go.

		if len(reqQuery.Data) == 0 {
			resQuery.Log = "Query cannot be zero length"
			resQuery.Code = wrsp.CodeType_EncodingError
			return
		}

		// set the query response height to current
		tree := app.state.Committed()

		height := reqQuery.Height
		if height == 0 {
			// TODO: once the rpc actually passes in non-zero
			// heights we can use to query right after a tx
			// we must retrun most recent, even if apphash
			// is not yet in the blockchain

			withProof := app.CommittedHeight() - 1
			if tree.Tree.VersionExists(withProof) {
				height = withProof
			} else {
				height = app.CommittedHeight()
			}
		}
		resQuery.Height = height

		switch reqQuery.Path {
		case "/store", "/key": // Get by key
			key := reqQuery.Data // Data holds the key bytes
			resQuery.Key = key
			if reqQuery.Prove {
				value, proof, err := tree.GetVersionedWithProof(key, height)
				if err != nil {
					resQuery.Log = err.Error()
					break
				}
				resQuery.Value = value
				resQuery.Proof = proof.Bytes()
			} else {
				value := tree.Get(key)
				resQuery.Value = value
			}

		default:
			resQuery.Code = wrsp.CodeType_UnknownRequest
			resQuery.Log = cmn.Fmt("Unexpected Query path: %v", reqQuery.Path)
		}
		return
	*/
}

// Commit implements wrsp.Application
func (app *App) Commit() (res wrsp.Result) {
	commitID := app.store.Commit()
	app.logger.Debug("Commit synced",
		"commit", commitID,
	)
	return wrsp.NewResultOK(hash, "")
}

// InitChain - WRSP
func (app *App) InitChain(req wrsp.RequestInitChain) {}

// BeginBlock - WRSP
func (app *App) BeginBlock(req wrsp.RequestBeginBlock) {
	app.header = req.Header
}

// EndBlock - WRSP
// Returns a list of all validator changes made in this block
func (app *App) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// XXX Update to res.Updates.
	res.Diffs = app.valSetDiff
	app.valSetDiff = nil
	return
}

// Return index of list with validator of same PubKey, or -1 if no match
func pubKeyIndex(val *wrsp.Validator, list []*wrsp.Validator) int {
	for i, v := range list {
		if bytes.Equal(val.PubKey, v.PubKey) {
			return i
		}
	}
	return -1
}

// Make a simple default logger
// TODO: Make log capturable for each transaction, and return it in
// ResponseDeliverTx.Log and ResponseCheckTx.Log.
func makeDefaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}

// InitState - used to setup state (was SetOption)
// to be call from setting up the genesis file
func (app *InitApp) InitState(module, key, value string) error {
	state := app.Append()
	logger := app.Logger().With("module", module, "key", key)

	if module == sdk.ModuleNameBase {
		if key == sdk.ChainKey {
			app.info.SetChainID(state, value)
			return nil
		}
		logger.Error("Invalid genesis option")
		return fmt.Errorf("Unknown base option: %s", key)
	}

	log, err := app.initState.InitState(logger, state, module, key, value)
	if err != nil {
		logger.Error("Invalid genesis option", "err", err)
	} else {
		logger.Info(log)
	}
	return err
}

// InitChain - WRSP - sets the initial validators
func (app *InitApp) InitChain(req wrsp.RequestInitChain) {
	// return early if no InitValidator registered
	if app.initVals == nil {
		return
	}

	logger, store := app.Logger(), app.Append()
	app.initVals.InitValidators(logger, store, req.Validators)
}
