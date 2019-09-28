package app

import (
	"bytes"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/types"
)

const mainKeyHeader = "header"

// App - The WRSP application
type App struct {
	logger log.Logger

	// App name from wrsp.Info
	name string

	// DeliverTx (main) state
	store types.MultiStore

	// CheckTx state
	storeCheck types.CacheMultiStore

	// Current block header
	header *wrsp.Header

	// Handler for CheckTx and DeliverTx.
	handler types.Handler

	// Cached validator changes from DeliverTx
	valUpdates []wrsp.Validator
}

var _ wrsp.Application = &App{}

func NewApp(name string) *App {
	return &App{
		name:   name,
		logger: makeDefaultLogger(),
	}
}

func (app *App) SetStore(store types.MultiStore) {
	app.store = store
}

func (app *App) SetHandler(handler types.Handler) {
	app.handler = handler
}

func (app *App) LoadLatestVersion() error {
	store := app.store
	store.LoadLatestVersion()
	return app.initFromStore()
}

func (app *App) LoadVersion(version int64) error {
	store := app.store
	store.LoadVersion(version)
	return app.initFromStore()
}

// Initializes the remaining logic from app.store.
func (app *App) initFromStore() error {
	store := app.store
	lastCommitID := store.LastCommitID()
	main := store.GetKVStore("main")
	header := (*wrsp.Header)(nil)
	storeCheck := store.CacheMultiStore()

	// Main store should exist.
	if store.GetKVStore("main") == nil {
		return errors.New("App expects MultiStore with 'main' KVStore")
	}

	// If we've committed before, we expect main://<mainKeyHeader>.
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
		if header.Height != lastCommitID.Version {
			errStr := fmt.Sprintf("Expected main://%s.Height %v but got %v", mainKeyHeader, version, headerHeight)
			return errors.New(errStr)
		}
	}

	// Set App state.
	app.header = header
	app.storeCheck = app.store.CacheMultiStore()
	app.valUpdates = nil

	return nil
}

//----------------------------------------

// Implements WRSP
func (app *App) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.store.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.Name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// Implements WRSP
func (app *App) SetOption(req wrsp.RequestSetOption) wrsp.ResponseSetOption {
	return "Not Implemented"
}

// Implements WRSP
func (app *App) InitChain(req wrsp.RequestInitChain) wrsp.ResponseInitChain {
	// TODO: Use req.Validators
}

// Implements WRSP
func (app *App) Query(req wrsp.RequestQuery) wrsp.ResponseQuery {
	// TODO: See app/query.go
}

// Implements WRSP
func (app *App) BeginBlock(req wrsp.RequestBeginBlock) wrsp.ResponseBeginBlock {
	app.header = req.Header
}

// Implements WRSP
func (app *App) CheckTx(txBytes []byte) wrsp.ResponseCheckTx {

	// Initialize arguments to Handler.
	var isCheckTx = true
	var ctx = types.NewContext(app.header, isCheckTx, txBytes)
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

// Implements WRSP
func (app *App) DeliverTx(txBytes []byte) wrsp.ResponseDeliverTx {

	// Initialize arguments to Handler.
	var isCheckTx = false
	var ctx = types.NewContext(app.header, isCheckTx, txBytes)
	var store = app.store
	var tx Tx = nil // nil until a decorator parses one.

	// Run the handler.
	var result = app.handler(ctx, app.store, tx)

	// After-handler hooks.
	if result.Code == wrsp.CodeType_OK {
		app.valUpdates = append(app.valUpdates, result.ValUpdate)
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

// Implements WRSP
func (app *App) EndBlock(req wrsp.RequestEndBlock) (res wrsp.ResponseEndBlock) {
	// XXX Update to res.Updates.
	res.ValidatorUpdates = app.valUpdates
	app.valUpdates = nil
	return
}

// Implements WRSP
func (app *App) Commit() (res wrsp.ResponseCommit) {
	commitID := app.store.Commit()
	app.logger.Debug("Commit synced",
		"commit", commitID,
	)
	return wrsp.NewResultOK(hash, "")
}

//----------------------------------------
// Misc.

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
