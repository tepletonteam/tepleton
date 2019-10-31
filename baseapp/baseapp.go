package baseapp

import (
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"

	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
)

// Key to store the header in the DB itself.
// Use the db directly instead of a store to avoid
// conflicts with handlers writing to the store
// and to avoid affecting the Merkle root.
var dbHeaderKey = []byte("header")

// The WRSP application
type BaseApp struct {
	// initialized on creation
	Logger     log.Logger
	name       string               // application name from wrsp.Info
	db         dbm.DB               // common DB backend
	cms        sdk.CommitMultiStore // Main (uncached) state
	router     Router               // handle any kind of message
	codespacer *sdk.Codespacer      // handle module codespacing

	// must be set
	txDecoder   sdk.TxDecoder   // unmarshal []byte into sdk.Tx
	anteHandler sdk.AnteHandler // ante handler for fee and auth

	// may be nil
	initChainer  sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes

	//--------------------
	// Volatile
	// checkState is set on initialization and reset on Commit.
	// deliverState is set in InitChain and BeginBlock and cleared on Commit.
	// See methods setCheckState and setDeliverState.
	// .valUpdates accumulate in DeliverTx and are reset in BeginBlock.
	// QUESTION: should we put valUpdates in the deliverState.ctx?
	checkState   *state           // for CheckTx
	deliverState *state           // for DeliverTx
	valUpdates   []wrsp.Validator // cached validator changes from DeliverTx
}

var _ wrsp.Application = (*BaseApp)(nil)

// Create and name new BaseApp
// NOTE: The db is used to store the version number for now.
func NewBaseApp(name string, logger log.Logger, db dbm.DB) *BaseApp {
	app := &BaseApp{
		Logger:     logger,
		name:       name,
		db:         db,
		cms:        store.NewCommitMultiStore(db),
		router:     NewRouter(),
		codespacer: sdk.NewCodespacer(),
	}
	// Register the undefined & root codespaces, which should not be used by any modules
	app.codespacer.RegisterOrPanic(sdk.CodespaceUndefined)
	app.codespacer.RegisterOrPanic(sdk.CodespaceRoot)
	return app
}

// BaseApp Name
func (app *BaseApp) Name() string {
	return app.name
}

// Register the next available codespace through the baseapp's codespacer, starting from a default
func (app *BaseApp) RegisterCodespace(codespace sdk.CodespaceType) sdk.CodespaceType {
	return app.codespacer.RegisterNext(codespace)
}

// Mount a store to the provided key in the BaseApp multistore
func (app *BaseApp) MountStoresIAVL(keys ...*sdk.KVStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeIAVL)
	}
}

// Mount a store to the provided key in the BaseApp multistore, using a specified DB
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// Mount a store to the provided key in the BaseApp multistore, using the default DB
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// nolint - Set functions
func (app *BaseApp) SetTxDecoder(txDecoder sdk.TxDecoder) {
	app.txDecoder = txDecoder
}
func (app *BaseApp) SetInitChainer(initChainer sdk.InitChainer) {
	app.initChainer = initChainer
}
func (app *BaseApp) SetBeginBlocker(beginBlocker sdk.BeginBlocker) {
	app.beginBlocker = beginBlocker
}
func (app *BaseApp) SetEndBlocker(endBlocker sdk.EndBlocker) {
	app.endBlocker = endBlocker
}
func (app *BaseApp) SetAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}

func (app *BaseApp) Router() Router { return app.router }

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

	// main store should exist.
	// TODO: we don't actually need the main store here
	main := app.cms.GetKVStore(mainKey)
	if main == nil {
		return errors.New("BaseApp expects MultiStore with 'main' KVStore")
	}

	// XXX: Do we really need the header? What does it have that we want
	// here that's not already in the CommitID ? If an app wants to have it,
	// they can do so in their BeginBlocker. If we force it in baseapp,
	// then either we force the AppHash to change with every block (since the header
	// will be in the merkle store) or we can't write the state and the header to the
	// db atomically without doing some surgery on the store interfaces ...

	// if we've committed before, we expect <dbHeaderKey> to exist in the db
	/*
		var lastCommitID = app.cms.LastCommitID()
		var header wrsp.Header

		if !lastCommitID.IsZero() {
			headerBytes := app.db.Get(dbHeaderKey)
			if len(headerBytes) == 0 {
				errStr := fmt.Sprintf("Version > 0 but missing key %s", dbHeaderKey)
				return errors.New(errStr)
			}
			err := proto.Unmarshal(headerBytes, &header)
			if err != nil {
				return errors.Wrap(err, "Failed to parse Header")
			}
			lastVersion := lastCommitID.Version
			if header.Height != lastVersion {
				errStr := fmt.Sprintf("Expected db://%s.Height %v but got %v", dbHeaderKey, lastVersion, header.Height)
				return errors.New(errStr)
			}
		}
	*/

	// initialize Check state
	app.setCheckState(wrsp.Header{})

	return nil
}

// NewContext returns a new Context with the correct store, the given header, and nil txBytes.
func (app *BaseApp) NewContext(isCheckTx bool, header wrsp.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, nil)
	}
	return sdk.NewContext(app.deliverState.ms, header, false, nil)
}

type state struct {
	ms  sdk.CacheMultiStore
	ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (app *BaseApp) setCheckState(header wrsp.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, true, nil),
	}
}

func (app *BaseApp) setDeliverState(header wrsp.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, nil),
	}
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

	// Initialize the deliver state and run initChain
	app.setDeliverState(wrsp.Header{})
	app.initChainer(app.deliverState.ctx, req) // no error

	// NOTE: we don't commit, but BeginBlock for block 1
	// starts from this deliverState

	return
}

// Implements WRSP.
// Delegates to CommitMultiStore if it implements Queryable
func (app *BaseApp) Query(req wrsp.RequestQuery) (res wrsp.ResponseQuery) {
	queryable, ok := app.cms.(sdk.Queryable)
	if !ok {
		msg := "application doesn't support queries"
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}
	return queryable.Query(req)
}

// Implements WRSP
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) (res wrsp.ResponseBeginBlock) {
	// Initialize the DeliverTx state.
	// If this is the first block, it should already
	// be initialized in InitChain. It may also be nil
	// if this is a test and InitChain was never called.
	if app.deliverState == nil {
		app.setDeliverState(req.Header)
	}
	app.valUpdates = nil
	if app.beginBlocker != nil {
		res = app.beginBlocker(app.deliverState.ctx, req)
	}
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
		// Even though the Result.Code is not OK, there are still effects,
		// namely fee deductions and sequence incrementing.
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

// Mostly for testing
func (app *BaseApp) Check(tx sdk.Tx) (result sdk.Result) {
	return app.runTx(true, nil, tx)
}
func (app *BaseApp) Deliver(tx sdk.Tx) (result sdk.Result) {
	return app.runTx(false, nil, tx)
}

// txBytes may be nil in some cases, eg. in tests.
// Also, in the future we may support "internal" transactions.
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
		err = err.WithDefaultCodespace(sdk.CodespaceRoot)
		return err.Result()
	}

	// Get the context
	var ctx sdk.Context
	if isCheckTx {
		ctx = app.checkState.ctx.WithTxBytes(txBytes)
	} else {
		ctx = app.deliverState.ctx.WithTxBytes(txBytes)
	}

	// Run the ante handler.
	if app.anteHandler != nil {
		newCtx, result, abort := app.anteHandler(ctx, tx)
		if abort {
			return result
		}
		if !newCtx.IsZero() {
			ctx = newCtx
		}
	}

	// Match route.
	msgType := msg.Type()
	handler := app.router.Route(msgType)
	if handler == nil {
		return sdk.ErrUnknownRequest("Unrecognized Msg type: " + msgType).Result()
	}

	// Get the correct cache
	var msCache sdk.CacheMultiStore
	if isCheckTx == true {
		// CacheWrap app.checkState.ms in case it fails.
		msCache = app.checkState.CacheMultiStore()
		ctx = ctx.WithMultiStore(msCache)
	} else {
		// CacheWrap app.deliverState.ms in case it fails.
		msCache = app.deliverState.CacheMultiStore()
		ctx = ctx.WithMultiStore(msCache)

	}

	result = handler(ctx, msg)

	// If result was successful, write to app.checkState.ms or app.deliverState.ms
	if result.IsOK() {
		msCache.Write()
	}

	return result
}

// Implements WRSP
func (app *BaseApp) EndBlock(req wrsp.RequestEndBlock) (res wrsp.ResponseEndBlock) {
	if app.endBlocker != nil {
		res = app.endBlocker(app.deliverState.ctx, req)
	} else {
		res.ValidatorUpdates = app.valUpdates
	}
	return
}

// Implements WRSP
func (app *BaseApp) Commit() (res wrsp.ResponseCommit) {
	header := app.deliverState.ctx.BlockHeader()
	/*
		// Write the latest Header to the store
			headerBytes, err := proto.Marshal(&header)
			if err != nil {
				panic(err)
			}
			app.db.SetSync(dbHeaderKey, headerBytes)
	*/

	// Write the Deliver state and commit the MultiStore
	app.deliverState.ms.Write()
	commitID := app.cms.Commit()
	app.Logger.Debug("Commit synced",
		"commit", commitID,
	)

	// Reset the Check state to the latest committed
	// NOTE: safe because Tendermint holds a lock on the mempool for Commit.
	// Use the header from this latest block.
	app.setCheckState(header)

	// Empty the Deliver state
	app.deliverState = nil

	return wrsp.ResponseCommit{
		Data: commitID.Hash,
	}
}
