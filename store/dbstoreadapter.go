package store

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	dbm "github.com/tepleton/tepleton/libs/db"
)

type dbStoreAdapter struct {
	dbm.DB
}

// Implements Store.
func (dbStoreAdapter) GetStoreType() StoreType {
	return sdk.StoreTypeDB
}

// Implements KVStore.
func (dsa dbStoreAdapter) CacheWrap() CacheWrap {
	return NewCacheKVStore(dsa)
}

// Implements KVStore
func (dsa dbStoreAdapter) Prefix(prefix []byte) KVStore {
	return prefixStore{dsa, prefix}
}

// dbm.DB implements KVStore so we can CacheKVStore it.
var _ KVStore = dbStoreAdapter{dbm.DB(nil)}
