package store

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	dbm "github.com/tepleton/tmlibs/db"
)

type dbStoreAdapter struct {
	dbm.DB
}

// Implements Store.
func (_ dbStoreAdapter) GetStoreType() StoreType {
	return sdk.StoreTypeDB
}

// Implements KVStore.
func (dsa dbStoreAdapter) CacheWrap() CacheWrap {
	return NewCacheKVStore(dsa)
}

// dbm.DB implements KVStore so we can CacheKVStore it.
var _ KVStore = dbStoreAdapter{dbm.DB(nil)}
