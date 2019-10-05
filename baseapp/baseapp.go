package baseapp

import (
	"fmt"
	"runtime/debug"

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

// The WRSP application
type BaseApp struct {
	logger      log.Logger
	name        string               // application name from wrsp.Info
	db          dbm.DB               // common DB backend
	cms         sdk.CommitMultiStore // Main (uncached) state
	txDecoder   sdk.TxDecoder        // unmarshal []byte into sdk.Tx
	initChainer sdk.InitChainer      //
	anteHandler sdk.AnteHandler      // ante handler for fee and auth
	router      Router               // handle any kind of message

	//--------------------
	// Volatile

	msCheck    sdk.CacheMultiStore // CheckTx state, a cache-wrap of `.cms`
	msDeliver  sdk.CacheMultiStore // DeliverTx state, a cache-wrap of `.cms`
	header     *wrsp.Header        // current block header
	valUpdates []wrsp.Validator    // cached validator changes from DeliverTx
}

var _ wrsp.Application = &BaseApp{}

// Create and name new BaseApp
func NewBaseApp(name string, logger log.Logger, db dbm.DB) *BaseApp {
	return &BaseApp{
		logger: logger,
		name:   name,
		db:     db,
		cms:    store.NewCommitMultiStore(db),
		router: NewRouter(),
	}
}

// BaseApp Name
func (app *BaseApp) Name() string {
	return app.name
}

// Mount a store to the provided key in the BaseApp multistore
func (app *BaseApp) MountStoresIAVL(keys ...*sdk.KVStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeIAVL)
	}
}

// Mount a store to the provided key in the BaseApp multistore
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, app.db)
}

// nolint - Set functions
func (app *BaseApp) SetTxDecoder(txDecoder sdk.TxDecoder) {
	app.txDecoder = txDecoder
}
func (app *BaseApp) SetInitChainer(initChainer sdk.InitChainer) {
	app.initChainer = initChainer
}
func (app *BaseApp) SetAnteHandler(ah sdk.AnteHandler) {
	// deducts fee from payer, verifies signatures and nonces, sets Signers to ctx.
	app.anteHandler = ah
}

// nolint - Get functions
func (app *BaseApp) Router() Router { return app.router }

/* TODO consider:
func (app *BaseApp) SetBeginBlocker(...) {}
func (app *BaseApp) SetEndBlocker(...) {}
*/

// load latest application version
func (app *BaseApp) LoadLatestVersion(mainKey sdk.StoreKey) error {
	app.cms.LoadLatestVersion()
	return app.initFromStore(mainKey)
}

// load application version
func (app *BaseApp) LoadVersion(version int64, mainKey sdk.StoreKey) error {
	app.cms.LoadVersion(version)
	return app.initFromStore(mainKey)
}

// the last CommitID of the multistore
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// the last commited block height
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromStore(mainKey sdk.StoreKey) error {
	var lastCommitID = app.cms.LastCommitID()
	var main = app.cms.GetKVStore(mainKey)
	var header *wrsp.Header

	// main store should exist.
	if main == nil {
		return errors.New("BaseApp expects MultiStore with 'main' KVStore")
	}

	// if we've committed before, we expect main://<mainHeaderKey>
	if !lastCommitID.IsZero() {
		headerBytes := main.Get(mainHeaderKey)
		if len(headerBytes) == 0 {
			errStr := fmt.Sprintf("Version > 0 but missing key %s", mainHeaderKey)
			return errors.New(errStr)
		}
		err := proto.Unmarshal(headerBytes, header)
		if err != nil {
			return errors.Wrap(err, "Failed to parse Header")
		}
		lastVersion := lastCommitID.Version
		if header.Height != lastVersion {
			errStr := fmt.Sprintf("Expected main://%s.Height %v but got %v", mainHeaderKey, lastVersion, header.Height)
			return errors.New(errStr)
		}
	}

	// set BaseApp state
	app.header = header
	app.msCheck = nil
	app.msDeliver = nil
	app.valUpdates = nil

	return nil
}

// NewContext returns a new Context suitable for AnteHandler (and indirectly Handler) processing.
// NOTE: txBytes may be nil to support TestApp.RunCheckTx
// and TestApp.RunDeliverTx.
func (app *BaseApp) NewContext(isCheckTx bool, txBytes []byte) sdk.Context {

	store := app.getMultiStore(isCheckTx)
	if store == nil {
		panic("BaseApp.NewContext() requires BeginBlock(): missing store")
	}
	if app.header == nil {
		panic("BaseApp.NewContext() requires BeginBlock(): missing header")
	}

	return sdk.NewContext(store, *app.header, isCheckTx, txBytes)
}

//----------------------------------------
// WRSP

// Implements WRSP
func (app *BaseApp) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {

	lastCommitID := app.cms.LastCommitID()

	return wrsp.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// Implements WRSP
func (app *BaseApp) SetOption(req wrsp.RequestSetOption) (res wrsp.ResponseSetOption) {
	// TODO: Implement
	return
}

// Implements WRSP
// InitChain runs the initialization logic directly on the CommitMultiStore and commits it.
func (app *BaseApp) InitChain(req wrsp.RequestInitChain) (res wrsp.ResponseInitChain) {
	if app.initChainer == nil {
		// TODO: should we have some default handling of validators?
		return
	}

	// make a context for the initialization.
	// NOTE: we're writing to the cms directly, without a CacheWrap
	ctx := sdk.NewContext(app.cms, wrsp.Header{}, false, nil)

	err := app.initChainer(ctx, req)
	if err != nil {
		// TODO: something better https://github.com/tepleton/tepleton-sdk/issues/468
		cmn.Exit(fmt.Sprintf("error initializing application genesis state: %v", err))
	}

	// XXX this commits everything and bumps the version.
	// With this, block 1 executes against state with version 1, but results in state with version 2.
	// Is that what we want ?
	app.cms.Commit()

	return
}

// Implements WRSP.
// Delegates to CommitMultiStore if it implements Queryable
func (app *BaseApp) Query(req wrsp.RequestQuery) (res wrsp.ResponseQuery) {
	queryable, ok := app.cms.(sdk.Queryable)
	if !ok {
		msg := "application doesn't support queries"
		return sdk.ErrUnknownRequest(msg).Result().ToQuery()
	}
	return queryable.Query(req)
}

// Implements WRSP
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) (res wrsp.ResponseBeginBlock) {
	// NOTE: For consistency we should unset these upon EndBlock.
	app.header = &req.Header
	app.msDeliver = app.cms.CacheMultiStore()
	app.msCheck = app.cms.CacheMultiStore()
	app.valUpdates = nil
	return
}

// Implements WRSP
func (app *BaseApp) CheckTx(txBytes []byte) (res wrsp.ResponseCheckTx) {

	// Decode the Tx.
	var result sdk.Result
	var tx, err = app.txDecoder(txBytes)
	if err != nil {
		result = err.Result()
	} else {
		result = app.runTx(true, txBytes, tx)
	}

	return wrsp.ResponseCheckTx{
		Code:      uint32(result.Code),
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
func (app *BaseApp) DeliverTx(txBytes []byte) (res wrsp.ResponseDeliverTx) {

	// Decode the Tx.
	var result sdk.Result
	var tx, err = app.txDecoder(txBytes)
	if err != nil {
		result = err.Result()
	} else {
		result = app.runTx(false, txBytes, tx)
	}

	// After-handler hooks.
	if result.IsOK() {
		app.valUpdates = append(app.valUpdates, result.ValidatorUpdates...)
	} else {
		// Even though the Code is not OK, there will be some side
		// effects, like those caused by fee deductions or sequence
		// incrementations.
	}

	// Tell the blockchain engine (i.e. Tendermint).
	return wrsp.ResponseDeliverTx{
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: result.GasWanted,
		GasUsed:   result.GasUsed,
		Tags:      result.Tags,
	}
}

// txBytes may be nil in some cases, for example, when tx is
// coming from TestApp.  Also, in the future we may support
// "internal" transactions.
func (app *BaseApp) runTx(isCheckTx bool, txBytes []byte, tx sdk.Tx) (result sdk.Result) {

	// Handle any panics.
	defer func() {
		if r := recover(); r != nil {
			log := fmt.Sprintf("Recovered: %v\nstack:\n%v", r, string(debug.Stack()))
			result = sdk.ErrInternal(log).Result()
		}
	}()

	// Get the Msg.
	var msg = tx.GetMsg()
	if msg == nil {
		return sdk.ErrInternal("Tx.GetMsg() returned nil").Result()
	}

	// Validate the Msg.
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	// Construct a Context.
	var ctx = app.NewContext(isCheckTx, txBytes)

	// TODO: override default ante handler w/ custom ante handler.

	// Run the ante handler.
	newCtx, result, abort := app.anteHandler(ctx, tx)
	if isCheckTx || abort {
		return result
	}
	if !newCtx.IsZero() {
		ctx = newCtx
	}

	// CacheWrap app.msDeliver in case it fails.
	msCache := app.getMultiStore(isCheckTx).CacheMultiStore()
	ctx = ctx.WithMultiStore(msCache)

	// Match and run route.
	msgType := msg.Type()
	handler := app.router.Route(msgType)
	result = handler(ctx, msg)

	// If result was successful, write to app.msDeliver or app.msCheck.
	if result.IsOK() {
		msCache.Write()
	}

	return result
}

// Implements WRSP
func (app *BaseApp) EndBlock(req wrsp.RequestEndBlock) (res wrsp.ResponseEndBlock) {
	res.ValidatorUpdates = app.valUpdates
	app.valUpdates = nil
	app.header = nil
	app.msDeliver = nil
	app.msCheck = nil
	return
}

// Implements WRSP
func (app *BaseApp) Commit() (res wrsp.ResponseCommit) {
	app.msDeliver.Write()
	commitID := app.cms.Commit()
	app.logger.Debug("Commit synced",
		"commit", commitID,
	)
	return wrsp.ResponseCommit{
		Data: commitID.Hash,
	}
}

//----------------------------------------
// Helpers

func (app *BaseApp) getMultiStore(isCheckTx bool) sdk.MultiStore {
	if isCheckTx {
		return app.msCheck
	}
	return app.msDeliver
}
