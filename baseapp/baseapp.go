package baseapp

import (
	"bytes"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
)

var mainHeaderKey = []byte("header")

// BaseApp - The WRSP application
type BaseApp struct {
	logger log.Logger

	// Application name from wrsp.Info
	name string

	// Common DB backend
	db dbm.DB

	// Main (uncached) state
	ms sdk.CommitMultiStore

	// Unmarshal []byte into sdk.Tx
	txDecoder sdk.TxDecoder

	// Ante handler for fee and auth.
	defaultAnteHandler sdk.AnteHandler

	// Handle any kind of message.
	router Router

	//--------------------
	// Volatile

	// CheckTx state, a cache-wrap of `.ms`.
	msCheck sdk.CacheMultiStore

	// DeliverTx state, a cache-wrap of `.ms`.
	msDeliver sdk.CacheMultiStore

	// Current block header
	header wrsp.Header

	// Cached validator changes from DeliverTx.
	valUpdates []wrsp.Validator
}

var _ wrsp.Application = &BaseApp{}

func NewBaseApp(name string) *BaseApp {
	var baseapp = &BaseApp{
		logger: makeDefaultLogger(),
		name:   name,
		db:     nil,
		ms:     nil,
		router: NewRouter(),
	}
	baseapp.initDB()
	baseapp.initMultiStore()
	return baseapp
}

// Create the underlying leveldb datastore which will
// persist the Merkle tree inner & leaf nodes.
func (app *BaseApp) initDB() {
	db, err := dbm.NewGoLevelDB(app.name, "data")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.db = db
}

func (app *BaseApp) initMultiStore() {
	ms := store.NewCommitMultiStore(app.db)
	app.ms = ms
}

func (app *BaseApp) Name() string {
	return app.name
}

func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.ms.MountStoreWithDB(key, typ, app.db)
}

func (app *BaseApp) TxDecoder() sdk.TxDecoder {
	return app.txDecoder
}

func (app *BaseApp) SetTxDecoder(txDecoder sdk.TxDecoder) {
	app.txDecoder = txDecoder
}

func (app *BaseApp) DefaultAnteHandler() sdk.AnteHandler {
	return app.defaultAnteHandler
}

func (app *BaseApp) SetDefaultAnteHandler(ah sdk.AnteHandler) {
	app.defaultAnteHandler = ah
}

func (app *BaseApp) MultiStore() sdk.MultiStore {
	return app.ms
}

func (app *BaseApp) MultiStoreCheck() sdk.MultiStore {
	return app.msCheck
}

func (app *BaseApp) MultiStoreDeliver() sdk.MultiStore {
	return app.msDeliver
}

func (app *BaseApp) Router() Router {
	return app.router
}

/* TODO consider:
func (app *BaseApp) SetBeginBlocker(...) {}
func (app *BaseApp) SetEndBlocker(...) {}
func (app *BaseApp) SetInitStater(...) {}
*/

func (app *BaseApp) LoadLatestVersion(mainKey sdk.StoreKey) error {
	app.ms.LoadLatestVersion()
	return app.initFromStore(mainKey)
}

func (app *BaseApp) LoadVersion(version int64, mainKey sdk.StoreKey) error {
	app.ms.LoadVersion(version)
	return app.initFromStore(mainKey)
}

// The last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.ms.LastCommitID()
}

// The last commited block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.ms.LastCommitID().Version
}

// Initializes the remaining logic from app.ms.
func (app *BaseApp) initFromStore(mainKey sdk.StoreKey) error {
	lastCommitID := app.ms.LastCommitID()
	main := app.ms.GetKVStore(mainKey)
	header := wrsp.Header{}

	// Main store should exist.
	if main == nil {
		return errors.New("BaseApp expects MultiStore with 'main' KVStore")
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

	// Set BaseApp state.
	app.header = header
	app.msCheck = nil
	app.msDeliver = nil
	app.valUpdates = nil

	return nil
}

//----------------------------------------

// Implements WRSP.
func (app *BaseApp) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.ms.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// Implements WRSP.
func (app *BaseApp) SetOption(req wrsp.RequestSetOption) (res wrsp.ResponseSetOption) {
	// TODO: Implement
	return
}

// Implements WRSP.
func (app *BaseApp) InitChain(req wrsp.RequestInitChain) (res wrsp.ResponseInitChain) {
	// TODO: Use req.Validators
	return
}

// Implements WRSP.
func (app *BaseApp) Query(req wrsp.RequestQuery) (res wrsp.ResponseQuery) {
	// TODO: See app/query.go
	return
}

// Implements WRSP.
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) (res wrsp.ResponseBeginBlock) {
	app.header = req.Header
	app.msDeliver = app.ms.CacheMultiStore()
	app.msCheck = app.ms.CacheMultiStore()
	return
}

// Implements WRSP.
func (app *BaseApp) CheckTx(txBytes []byte) (res wrsp.ResponseCheckTx) {

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

// Implements WRSP.
func (app *BaseApp) DeliverTx(txBytes []byte) (res wrsp.ResponseDeliverTx) {

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

func (app *BaseApp) runTx(isCheckTx bool, txBytes []byte) (result sdk.Result) {

	// Handle any panics.
	defer func() {
		if r := recover(); r != nil {
			result = sdk.Result{
				Code: 1, // TODO
				Log:  fmt.Sprintf("Recovered: %v\n", r),
			}
		}
	}()

	// Construct a Context.
	var ctx = app.NewContext(isCheckTx, txBytes)

	// Decode the Tx.
	tx, err := app.txDecoder(txBytes)
	if err != nil {
		return sdk.Result{
			Code: 1, //  TODO
		}
	}

	// TODO: override default ante handler w/ custom ante handler.

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

// Implements WRSP.
func (app *BaseApp) EndBlock(req wrsp.RequestEndBlock) (res wrsp.ResponseEndBlock) {
	res.ValidatorUpdates = app.valUpdates
	app.valUpdates = nil
	return
}

// Implements WRSP.
func (app *BaseApp) Commit() (res wrsp.ResponseCommit) {
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
