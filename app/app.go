package app

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/iavl"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/errors"
	sm "github.com/tepleton/tepleton-sdk/state"
	"github.com/tepleton/tepleton-sdk/version"
)

//nolint
const (
	ChainKey = "chain_id"
)

// BaseApp contains a data store and all info needed
// to perform queries and handshakes.
//
// It should be embeded in another struct for CheckTx,
// DeliverTx and initializing state from the genesis.
type BaseApp struct {
	info *sm.ChainState
	*sm.State

	// cached validator changes from DeliverTx
	pending []*wrsp.Validator

	// height is last committed block, DeliverTx is the next one
	height uint64

	logger log.Logger
}

// NewBaseApp creates a data store to handle queries
func NewBaseApp(dbName string, cacheSize int, logger log.Logger) (*BaseApp, error) {
	state, err := loadState(dbName, cacheSize)
	if err != nil {
		return nil, err
	}
	app := &BaseApp{
		State:  state,
		height: state.LatestHeight(),
		info:   sm.NewChainState(),
		logger: logger,
	}
	return app, nil
}

// GetChainID returns the currently stored chain
func (app *BaseApp) GetChainID() string {
	return app.info.GetChainID(app.Committed())
}

// Logger returns the application base logger
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// Hash gets the last hash stored in the database
func (app *BaseApp) Hash() []byte {
	return app.State.LatestHash()
}

// Info implements wrsp.Application. It returns the height and hash,
// as well as the wrsp name and version.
//
// The height is the block that holds the transactions, not the apphash itself.
func (app *BaseApp) Info(req wrsp.RequestInfo) wrsp.ResponseInfo {
	hash := app.Hash()

	app.logger.Info("Info synced",
		"height", app.height,
		"hash", fmt.Sprintf("%X", hash))

	return wrsp.ResponseInfo{
		// TODO
		Data:             fmt.Sprintf("Basecoin v%v", version.Version),
		LastBlockHeight:  app.height,
		LastBlockAppHash: hash,
	}
}

// SetOption - WRSP
func (app *BaseApp) SetOption(key string, value string) string {
	return "Not Implemented"
}

// Query - WRSP
func (app *BaseApp) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Log = "Query cannot be zero length"
		resQuery.Code = wrsp.CodeType_EncodingError
		return
	}

	// set the query response height to current
	tree := app.State.Committed()

	height := reqQuery.Height
	if height == 0 {
		// TODO: once the rpc actually passes in non-zero
		// heights we can use to query right after a tx
		// we must retrun most recent, even if apphash
		// is not yet in the blockchain

		// if tree.Tree.VersionExists(app.height - 1) {
		//  height = app.height - 1
		// } else {
		height = app.height
		// }
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
}

// Commit implements wrsp.Application
func (app *BaseApp) Commit() (res wrsp.Result) {
	app.height++

	hash, err := app.State.Commit(app.height)
	if err != nil {
		// die if we can't commit, not to recover
		panic(err)
	}
	app.logger.Debug("Commit synced",
		"height", app.height,
		"hash", fmt.Sprintf("%X", hash),
	)

	if app.State.Size() == 0 {
		return wrsp.NewResultOK(nil, "Empty hash for empty tree")
	}
	return wrsp.NewResultOK(hash, "")
}

// InitChain - WRSP
func (app *BaseApp) InitChain(req wrsp.RequestInitChain) {
}

// BeginBlock - WRSP
func (app *BaseApp) BeginBlock(req wrsp.RequestBeginBlock) {
}

// EndBlock - WRSP
// Returns a list of all validator changes made in this block
func (app *BaseApp) EndBlock(height uint64) (res wrsp.ResponseEndBlock) {
	// TODO: cleanup in case a validator exists multiple times in the list
	res.Diffs = app.pending
	app.pending = nil
	return
}

// AddValChange is meant to be called by apps on DeliverTx
// results, this is added to the cache for the endblock
// changeset
func (app *BaseApp) AddValChange(diffs []*wrsp.Validator) {
	for _, d := range diffs {
		idx := pubKeyIndex(d, app.pending)
		if idx >= 0 {
			app.pending[idx] = d
		} else {
			app.pending = append(app.pending, d)
		}
	}
}

// return index of list with validator of same PubKey, or -1 if no match
func pubKeyIndex(val *wrsp.Validator, list []*wrsp.Validator) int {
	for i, v := range list {
		if bytes.Equal(val.PubKey, v.PubKey) {
			return i
		}
	}
	return -1
}

func loadState(dbName string, cacheSize int) (*sm.State, error) {
	// memory backed case, just for testing
	if dbName == "" {
		tree := iavl.NewVersionedTree(0, dbm.NewMemDB())
		return sm.NewState(tree), nil
	}

	// Expand the path fully
	dbPath, err := filepath.Abs(dbName)
	if err != nil {
		return nil, errors.ErrInternal("Invalid Database Name")
	}

	// Some external calls accidently add a ".db", which is now removed
	dbPath = strings.TrimSuffix(dbPath, path.Ext(dbPath))

	// Split the database name into it's components (dir, name)
	dir := path.Dir(dbPath)
	name := path.Base(dbPath)

	// Make sure the path exists
	empty, _ := cmn.IsDirEmpty(dbPath + ".db")

	// Open database called "dir/name.db", if it doesn't exist it will be created
	db := dbm.NewDB(name, dbm.LevelDBBackendStr, dir)
	tree := iavl.NewVersionedTree(cacheSize, db)

	if !empty {
		if err = tree.Load(); err != nil {
			return nil, errors.ErrInternal("Loading tree: " + err.Error())
		}
	}

	return sm.NewState(tree), nil
}
