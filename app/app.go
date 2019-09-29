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

	"github.com/tepleton/tepleton-sdk/types"
)

var mainHeaderKey = []byte("header")

// App - The WRSP application
type App struct {
	logger log.Logger

	// App name from wrsp.Info
	name string

	// Main (uncached) state
	ms types.CommitMultiStore

	// Unmarshal []byte into types.Tx
	txDecoder types.TxDecoder

	// Ante handler for fee and auth.
	defaultAnteHandler types.AnteHandler

	// Handle any kind of message.
	router Router

	//--------------------
	// Volatile

	// CheckTx state, a cache-wrap of `.ms`.
	msCheck types.CacheMultiStore

	// DeliverTx state, a cache-wrap of `.ms`.
	msDeliver types.CacheMultiStore

	// Current block header
	header wrsp.Header

	// Cached validator changes from DeliverTx.
	valUpdates []wrsp.Validator
}

var _ wrsp.Application = &App{}

func NewApp(name string, ms CommitMultiStore) *App {
	return &App{
		logger: makeDefaultLogger(),
		name:   name,
		ms:     ms,
		router: NewRouter(),
	}
}

func (app *App) Name() string {
	return app.name
}

func (app *App) SetTxDecoder(txDecoder types.TxDecoder) {
	app.txDecoder = txDecoder
}

func (app *App) SetDefaultAnteHandler(ah types.AnteHandler) {
	app.defaultAnteHandler = ah
}

func (app *App) Router() Router {
	return app.router
}

/* TODO consider:
func (app *App) SetBeginBlocker(...) {}
func (app *App) SetEndBlocker(...) {}
func (app *App) SetInitStater(...) {}
*/

func (app *App) LoadLatestVersion() error {
	app.ms.LoadLatestVersion()
	return app.initFromStore()
}

func (app *App) LoadVersion(version int64) error {
	app.ms.LoadVersion(version)
	return app.initFromStore()
}

// The last CommitID of the multistore.
func (app *App) LastCommitID() types.CommitID {
	return app.ms.LastCommitID()
}

// The last commited block height.
func (app *App) LastBlockHeight() int64 {
	return app.ms.LastCommitID().Version
}

// Initializes the remaining logic from app.ms.
func (app *App) initFromStore() error {
	lastCommitID := app.ms.LastCommitID()
	main := app.ms.GetKVStore("main")
	header := wrsp.Header{}

	// Main store should exist.
	if app.ms.GetKVStore("main") == nil {
		return errors.New("App expects MultiStore with 'main' KVStore")
	}

	// If we've committed before, we expect main://<mainHeaderKey>.
	if !lastCommitID.IsZero() {
		headerBytes := main.Get(mainHeaderKey)
		if len(headerBytes) == 0 {
			errStr := fmt.Sprintf("Version > 0 but missing key %s", mainHeaderKey)
			return errors.New(errStr)
		}
		err := proto.Unmarshal(headerBytes, &header)
		if err != nil {
			return errors.Wrap(err, "Failed to parse Header")
		}
		lastVersion := lastCommitID.Version
		if header.Height != lastVersion {
			errStr := fmt.Sprintf("Expected main://%s.Height %v but got %v", mainHeaderKey, lastVersion, header.Height)
			return errors.New(errStr)
		}
	}

	// Set App state.
	app.header = header
	app.msCheck = nil
	app.msDeliver = nil
	app.valUpdates = nil

	return nil
}

//----------------------------------------

// Implements WRSP
func (app *App) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.ms.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// Implements WRSP
func (app *App) SetOption(req wrsp.RequestSetOption) (res wrsp.ResponseSetOption) {
	// TODO: Implement
	return
}

// Implements WRSP
func (app *App) InitChain(req wrsp.RequestInitChain) (res wrsp.ResponseInitChain) {
	// TODO: Use req.Validators
	return
}

// Implements WRSP
func (app *App) Query(req wrsp.RequestQuery) (res wrsp.ResponseQuery) {
	// TODO: See app/query.go
	return
}

// Implements WRSP
func (app *App) BeginBlock(req wrsp.RequestBeginBlock) (res wrsp.ResponseBeginBlock) {
	app.header = req.Header
	app.msDeliver = app.ms.CacheMultiStore()
	app.msCheck = app.ms.CacheMultiStore()
	return
}

// Implements WRSP
func (app *App) CheckTx(txBytes []byte) (res wrsp.ResponseCheckTx) {

	result := app.runTx(true, txBytes)

	return wrsp.ResponseCheckTx{
		Code:      result.Code,
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: result.GasWanted,
		Fee: cmn.KI64Pair{
			[]byte(result.FeeDenom),
			result.FeeAmount,
		},
		Tags: result.Tags,
	}

}

// Implements WRSP
func (app *App) DeliverTx(txBytes []byte) (res wrsp.ResponseDeliverTx) {

	result := app.runTx(false, txBytes)

	// After-handler hooks.
	if result.Code == wrsp.CodeTypeOK {
		app.valUpdates = append(app.valUpdates, result.ValidatorUpdates...)
	} else {
		// Even though the Code is not OK, there will be some side
		// effects, like those caused by fee deductions or sequence
		// incrementations.
	}

	// Tell the blockchain engine (i.e. Tendermint).
	return wrsp.ResponseDeliverTx{
		Code:      result.Code,
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: result.GasWanted,
		GasUsed:   result.GasUsed,
		Tags:      result.Tags,
	}
}

func (app *App) runTx(isCheckTx bool, txBytes []byte) (result types.Result) {

	// Handle any panics.
	defer func() {
		if r := recover(); r != nil {
			result = types.Result{
				Code: 1, // TODO
				Log:  fmt.Sprintf("Recovered: %v\n", r),
			}
		}
	}()

	var store types.MultiStore
	if isCheckTx {
		store = app.msCheck
	} else {
		store = app.msDeliver
	}

	// Initialize arguments to Handler.
	var ctx = types.NewContext(
		store,
		app.header,
		isCheckTx,
		txBytes,
	)

	// Decode the Tx.
	var err error
	tx, err = app.txDecoder(txBytes)
	if err != nil {
		return types.Result{
			Code: 1, //  TODO
		}
	}

	// Run the ante handler.
	ctx, result, abort := app.defaultAnteHandler(ctx, tx)
	if isCheckTx || abort {
		return result
	}

	// Match and run route.
	msgType := tx.Type()
	handler := app.router.Route(msgType)
	result = handler(ctx, tx)

	return result
}

// Implements WRSP
func (app *App) EndBlock(req wrsp.RequestEndBlock) (res wrsp.ResponseEndBlock) {
	res.ValidatorUpdates = app.valUpdates
	app.valUpdates = nil
	return
}

// Implements WRSP
func (app *App) Commit() (res wrsp.ResponseCommit) {
	app.msDeliver.Write()
	commitID := app.ms.Commit()
	app.logger.Debug("Commit synced",
		"commit", commitID,
	)
	return wrsp.ResponseCommit{
		Data: commitID.Hash,
	}
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
