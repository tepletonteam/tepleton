package app

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-wire"
	"github.com/tepleton/iavl"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/state"
)

// Store contains the merkle tree, and all info to handle wrsp requests
type Store struct {
	state.State
	height    uint64
	hash      []byte
	persisted bool

	logger log.Logger
}

var stateKey = []byte("merkle:state") // Database key for merkle tree save value db values

// ChainState contains the latest Merkle root hash and the number of times `Commit` has been called
type ChainState struct {
	Hash   []byte
	Height uint64
}

// MockStore returns an in-memory store only intended for testing
func MockStore() *Store {
	res, err := NewStore("", 0, log.NewNopLogger())
	if err != nil {
		// should never happen, abort test if it does
		panic(err)
	}
	return res
}

// NewStore initializes an in-memory IAVLTree, or attempts to load a persistant
// tree from disk
func NewStore(dbName string, cacheSize int, logger log.Logger) (*Store, error) {
	// start at 1 so the height returned by query is for the
	// next block, ie. the one that includes the AppHash for our current state
	initialHeight := uint64(1)

	// Non-persistent case
	if dbName == "" {
		tree := iavl.NewIAVLTree(
			0,
			nil,
		)
		store := &Store{
			State:  state.NewState(tree, false),
			height: initialHeight,
			logger: logger,
		}
		return store, nil
	}

	// Expand the path fully
	dbPath, err := filepath.Abs(dbName)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Database Name")
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
	tree := iavl.NewIAVLTree(cacheSize, db)

	var chainState ChainState
	if empty {
		logger.Info("no existing db, creating new db")
		chainState = ChainState{
			Hash:   tree.Save(),
			Height: initialHeight,
		}
		db.Set(stateKey, wire.BinaryBytes(chainState))
	} else {
		logger.Info("loading existing db")
		eyesStateBytes := db.Get(stateKey)
		err = wire.ReadBinaryBytes(eyesStateBytes, &chainState)
		if err != nil {
			return nil, errors.Wrap(err, "Reading MerkleEyesState")
		}
		tree.Load(chainState.Hash)
	}

	res := &Store{
		State:     state.NewState(tree, true),
		height:    chainState.Height,
		hash:      chainState.Hash,
		persisted: true,
		logger:    logger,
	}
	return res, nil
}

// Info implements wrsp.Application. It returns the height, hash and size (in the data).
// The height is the block that holds the transactions, not the apphash itself.
func (s *Store) Info() wrsp.ResponseInfo {
	s.logger.Info("Info synced",
		"height", s.height,
		"hash", fmt.Sprintf("%X", s.hash))
	return wrsp.ResponseInfo{
		Data:             cmn.Fmt("size:%v", s.State.Size()),
		LastBlockHeight:  s.height - 1,
		LastBlockAppHash: s.hash,
	}
}

// Commit implements wrsp.Application
func (s *Store) Commit() wrsp.Result {
	var err error
	s.height++
	s.hash, err = s.State.Hash()
	if err != nil {
		return wrsp.NewError(wrsp.CodeType_InternalError, err.Error())
	}

	s.logger.Debug("Commit synced",
		"height", s.height,
		"hash", fmt.Sprintf("%X", s.hash))

	s.State.BatchSet(stateKey, wire.BinaryBytes(ChainState{
		Hash:   s.hash,
		Height: s.height,
	}))

	hash, err := s.State.Commit()
	if err != nil {
		return wrsp.NewError(wrsp.CodeType_InternalError, err.Error())
	}
	if !bytes.Equal(hash, s.hash) {
		return wrsp.NewError(wrsp.CodeType_InternalError, "AppHash is incorrect")
	}

	if s.State.Size() == 0 {
		return wrsp.NewResultOK(nil, "Empty hash for empty tree")
	}
	return wrsp.NewResultOK(s.hash, "")
}

// Query implements wrsp.Application
func (s *Store) Query(reqQuery wrsp.RequestQuery) (resQuery wrsp.ResponseQuery) {

	if reqQuery.Height != 0 {
		// TODO: support older commits
		resQuery.Code = wrsp.CodeType_InternalError
		resQuery.Log = "merkleeyes only supports queries on latest commit"
		return
	}

	// set the query response height to current
	resQuery.Height = s.height

	tree := s.State.Committed()

	switch reqQuery.Path {
	case "/store", "/key": // Get by key
		key := reqQuery.Data // Data holds the key bytes
		resQuery.Key = key
		if reqQuery.Prove {
			value, proof, err := tree.GetWithProof(key)
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
