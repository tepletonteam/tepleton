package store

import (
	"fmt"
	"sync"

	"github.com/tepleton/go-amino"
	"github.com/tepleton/iavl"
	wrsp "github.com/tepleton/tepleton/wrsp/types"
	cmn "github.com/tepleton/tepleton/libs/common"
	dbm "github.com/tepleton/tepleton/libs/db"

	sdk "github.com/tepleton/tepleton-sdk/types"
)

const (
	defaultIAVLCacheSize  = 10000
	defaultIAVLNumRecent  = 100
	defaultIAVLStoreEvery = 10000
)

// load the iavl store
func LoadIAVLStore(db dbm.DB, id CommitID) (CommitStore, error) {
	tree := iavl.NewVersionedTree(db, defaultIAVLCacheSize)
	_, err := tree.LoadVersion(id.Version)
	if err != nil {
		return nil, err
	}
	store := newIAVLStore(tree, defaultIAVLNumRecent, defaultIAVLStoreEvery)
	return store, nil
}

//----------------------------------------

var _ KVStore = (*iavlStore)(nil)
var _ CommitStore = (*iavlStore)(nil)
var _ Queryable = (*iavlStore)(nil)

// iavlStore Implements KVStore and CommitStore.
type iavlStore struct {

	// The underlying tree.
	tree *iavl.VersionedTree

	// How many old versions we hold onto.
	// A value of 0 means keep no recent states
	numRecent int64

	// Distance between state-sync waypoint states to be stored
	// See https://github.com/tepleton/tepleton/issues/828
	// A value of 1 means store every state
	// A value of 0 means store no waypoints (node cannot assist in state-sync)
	// By default this value should be set the same across all nodes,
	// so that nodes can know the waypoints their peers store
	// TODO if set to non-default, signal to peers that the node is not suitable as a state sync source
	storeEvery int64
}

// CONTRACT: tree should be fully loaded.
func newIAVLStore(tree *iavl.VersionedTree, numRecent int64, storeEvery int64) *iavlStore {
	st := &iavlStore{
		tree:       tree,
		numRecent:  numRecent,
		storeEvery: storeEvery,
	}
	return st
}

// Implements Committer.
func (st *iavlStore) Commit() CommitID {

	// Save a new version.
	hash, version, err := st.tree.SaveVersion()
	if err != nil {
		// TODO: Do we want to extend Commit to allow returning errors?
		panic(err)
	}

	// Release an old version of history, if not a sync waypoint
	previous := version - 1
	if st.numRecent < previous {
		toRelease := previous - st.numRecent
		if st.storeEvery == 0 || toRelease%st.storeEvery != 0 {
			err := st.tree.DeleteVersion(toRelease)
			if err != nil {
				panic(err)
			}
		}
	}

	return CommitID{
		Version: version,
		Hash:    hash,
	}
}

// Implements Committer.
func (st *iavlStore) LastCommitID() CommitID {
	return CommitID{
		Version: st.tree.Version64(),
		Hash:    st.tree.Hash(),
	}
}

// VersionExists returns whether or not a given version is stored
func (st *iavlStore) VersionExists(version int64) bool {
	return st.tree.VersionExists(version)
}

// Implements Store.
func (st *iavlStore) GetStoreType() StoreType {
	return sdk.StoreTypeIAVL
}

// Implements Store.
func (st *iavlStore) CacheWrap() CacheWrap {
	return NewCacheKVStore(st)
}

// Implements KVStore.
func (st *iavlStore) Set(key, value []byte) {
	st.tree.Set(key, value)
}

// Implements KVStore.
func (st *iavlStore) Get(key []byte) (value []byte) {
	_, v := st.tree.Get(key)
	return v
}

// Implements KVStore.
func (st *iavlStore) Has(key []byte) (exists bool) {
	return st.tree.Has(key)
}

// Implements KVStore.
func (st *iavlStore) Delete(key []byte) {
	st.tree.Remove(key)
}

// Implements KVStore
func (st *iavlStore) Prefix(prefix []byte) KVStore {
	return prefixStore{st, prefix}
}

// Implements KVStore.
func (st *iavlStore) Iterator(start, end []byte) Iterator {
	return newIAVLIterator(st.tree.Tree(), start, end, true)
}

// Implements KVStore.
func (st *iavlStore) ReverseIterator(start, end []byte) Iterator {
	return newIAVLIterator(st.tree.Tree(), start, end, false)
}

// Query implements WRSP interface, allows queries
//
// by default we will return from (latest height -1),
// as we will have merkle proofs immediately (header height = data height + 1)
// If latest-1 is not present, use latest (which must be present)
// if you care to have the latest data to see a tx results, you must
// explicitly set the height you want to see
func (st *iavlStore) Query(req wrsp.RequestQuery) (res wrsp.ResponseQuery) {
	if len(req.Data) == 0 {
		msg := "Query cannot be zero length"
		return sdk.ErrTxDecode(msg).QueryResult()
	}

	tree := st.tree
	height := req.Height
	if height == 0 {
		latest := tree.Version64()
		if tree.VersionExists(latest - 1) {
			height = latest - 1
		} else {
			height = latest
		}
	}
	// store the height we chose in the response
	res.Height = height

	switch req.Path {
	case "/store", "/key": // Get by key
		key := req.Data // Data holds the key bytes
		res.Key = key
		if req.Prove {
			value, proof, err := tree.GetVersionedWithProof(key, height)
			if err != nil {
				res.Log = err.Error()
				break
			}
			res.Value = value
			cdc := amino.NewCodec()
			p, err := cdc.MarshalBinary(proof)
			if err != nil {
				res.Log = err.Error()
				break
			}
			res.Proof = p
		} else {
			_, res.Value = tree.GetVersioned(key, height)
		}
	case "/subspace":
		subspace := req.Data
		res.Key = subspace
		var KVs []KVPair
		iterator := sdk.KVStorePrefixIterator(st, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, KVPair{iterator.Key(), iterator.Value()})
		}
		iterator.Close()
		res.Value = cdc.MustMarshalBinary(KVs)
	default:
		msg := fmt.Sprintf("Unexpected Query path: %v", req.Path)
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}
	return
}

//----------------------------------------

// Implements Iterator.
type iavlIterator struct {
	// Underlying store
	tree *iavl.Tree

	// Domain
	start, end []byte

	// Iteration order
	ascending bool

	// Channel to push iteration values.
	iterCh chan cmn.KVPair

	// Close this to release goroutine.
	quitCh chan struct{}

	// Close this to signal that state is initialized.
	initCh chan struct{}

	//----------------------------------------
	// What follows are mutable state.
	mtx sync.Mutex

	invalid bool   // True once, true forever
	key     []byte // The current key
	value   []byte // The current value
}

var _ Iterator = (*iavlIterator)(nil)

// newIAVLIterator will create a new iavlIterator.
// CONTRACT: Caller must release the iavlIterator, as each one creates a new
// goroutine.
func newIAVLIterator(tree *iavl.Tree, start, end []byte, ascending bool) *iavlIterator {
	iter := &iavlIterator{
		tree:      tree,
		start:     cp(start),
		end:       cp(end),
		ascending: ascending,
		iterCh:    make(chan cmn.KVPair, 0), // Set capacity > 0?
		quitCh:    make(chan struct{}),
		initCh:    make(chan struct{}),
	}
	go iter.iterateRoutine()
	go iter.initRoutine()
	return iter
}

// Run this to funnel items from the tree to iterCh.
func (iter *iavlIterator) iterateRoutine() {
	iter.tree.IterateRange(
		iter.start, iter.end, iter.ascending,
		func(key, value []byte) bool {
			select {
			case <-iter.quitCh:
				return true // done with iteration.
			case iter.iterCh <- cmn.KVPair{key, value}:
				return false // yay.
			}
		},
	)
	close(iter.iterCh) // done.
}

// Run this to fetch the first item.
func (iter *iavlIterator) initRoutine() {
	iter.receiveNext()
	close(iter.initCh)
}

// Implements Iterator.
func (iter *iavlIterator) Domain() (start, end []byte) {
	return iter.start, iter.end
}

// Implements Iterator.
func (iter *iavlIterator) Valid() bool {
	iter.waitInit()
	iter.mtx.Lock()
	defer iter.mtx.Unlock()

	return !iter.invalid
}

// Implements Iterator.
func (iter *iavlIterator) Next() {
	iter.waitInit()
	iter.mtx.Lock()
	defer iter.mtx.Unlock()
	iter.assertIsValid()

	iter.receiveNext()
}

// Implements Iterator.
func (iter *iavlIterator) Key() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	defer iter.mtx.Unlock()
	iter.assertIsValid()

	return iter.key
}

// Implements Iterator.
func (iter *iavlIterator) Value() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	defer iter.mtx.Unlock()
	iter.assertIsValid()

	return iter.value
}

// Implements Iterator.
func (iter *iavlIterator) Close() {
	close(iter.quitCh)
}

//----------------------------------------

func (iter *iavlIterator) setNext(key, value []byte) {
	iter.assertIsValid()

	iter.key = key
	iter.value = value
}

func (iter *iavlIterator) setInvalid() {
	iter.assertIsValid()

	iter.invalid = true
}

func (iter *iavlIterator) waitInit() {
	<-iter.initCh
}

func (iter *iavlIterator) receiveNext() {
	kvPair, ok := <-iter.iterCh
	if ok {
		iter.setNext(kvPair.Key, kvPair.Value)
	} else {
		iter.setInvalid()
	}
}

func (iter *iavlIterator) assertIsValid() {
	if iter.invalid {
		panic("invalid iterator")
	}
}

//----------------------------------------

func cp(bz []byte) (ret []byte) {
	if bz == nil {
		return nil
	}
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}
