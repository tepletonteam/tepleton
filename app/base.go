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

// BaseApp - The WRSP application
type BaseApp struct {
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

var _ wrsp.Application = &BaseApp{}

// CONTRACT: There exists a "main" KVStore.
func NewBaseApp(name string, store CommitMultiStore) (*BaseApp, error) {

	if store.GetKVStore("main") == nil {
		return nil, errors.New("BaseApp expects MultiStore with 'main' KVStore")
	}

	logger := makeDefaultLogger()
	lastCommitID := store.LastCommitID()
	curVersion := store.CurrentVersion()
	main := store.GetKVStore("main")
	header := (*wrsp.Header)(nil)
	storeCheck := store.CacheMultiStore()

	// SANITY
	if curVersion != lastCommitID.Version+1 {
		panic("CurrentVersion != LastCommitID.Version+1")
	}

	// If we've committed before, we expect store.GetKVStore("main").Get("header")
	if !lastCommitID.IsZero() {
		headerBytes, ok := main.Get(mainKeyHeader)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Version > 0 but missing key %s", mainKeyHeader))
		}
		err = proto.Unmarshal(headerBytes, header)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse Header")
		}

		// SANITY: Validate Header
		if header.Height != curVersion-1 {
			errStr := fmt.Sprintf("Expected header.Height %v but got %v", version, headerHeight)
			panic(errStr)
		}
	}

	return &BaseApp{
		logger:     logger,
		name:       name,
		store:      store,
		storeCheck: storeCheck,
		header:     header,
		hander:     nil, // set w/ .WithHandler()
		valSetDiff: nil,
	}
}

func (app *BaseApp) WithHandler(handler sdk.Handler) *BaseApp {
	app.handler = handler
}

//----------------------------------------

// DeliverTx - WRSP - dispatches to the handler
func (app *BaseApp) DeliverTx(txBytes []byte) wrsp.ResponseDeliverTx {

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
func (app *BaseApp) CheckTx(txBytes []byte) wrsp.ResponseCheckTx {

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
func (app *BaseApp) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.store.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.Name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// SetOption - WRSP
func (app *BaseApp) SetOption(key string, value string) string {
	return "Not Implemented"
}

// Query - WRSP
func (app *BaseApp) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
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
func (app *BaseApp) Commit() (res wrsp.Result) {
	/*
		hash, err := app.state.Commit(app.height)
		if err != nil {
			// die if we can't commit, not to recover
			panic(err)
		}
		app.logger.Debug("Commit synced",
			"height", app.height,
			"hash", fmt.Sprintf("%X", hash),
		)

		if app.state.Size() == 0 {
			return wrsp.NewResultOK(nil, "Empty hash for empty tree")
		}
		return wrsp.NewResultOK(hash, "")
	*/
}

// InitChain - WRSP
func (app *BaseApp) InitChain(req wrsp.RequestInitChain) {}

// BeginBlock - WRSP
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) {
	app.header = req.Header
}

// EndBlock - WRSP
// Returns a list of all validator changes made in this block
func (app *BaseApp) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// TODO: Compress duplicates
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
