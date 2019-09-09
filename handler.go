package basecoin

import (
	"bytes"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-wire/data"

	"github.com/tepleton/basecoin/types"
)

// Handler is anything that processes a transaction
type Handler interface {
	CheckTx(ctx Context, store types.KVStore, tx Tx) (Result, error)
	DeliverTx(ctx Context, store types.KVStore, tx Tx) (Result, error)

	// TODO: flesh these out as well
	// SetOption(store types.KVStore, key, value string) (log string)
	// InitChain(store types.KVStore, vals []*wrsp.Validator)
	// BeginBlock(store types.KVStore, hash []byte, header *wrsp.Header)
	// EndBlock(store types.KVStore, height uint64) wrsp.ResponseEndBlock
}

// different apps to authorize
const (
	Sigs = "sigs"
	IBC  = "ibc"
	Role = "role"
)

// TODO: handle this in some secure way, only certain apps can add permissions
type Permission struct {
	App     string // Which app authorized this?
	Address []byte // App-specific identifier
}

// TODO: Context is a place-holder, soon we add some request data here from the
// higher-levels (like tell an app who signed).
// Trust me, we will need it like CallContext now...
type Context struct {
	perms []Permission
}

// TOTALLY insecure.  will redo later, but you get the point
func (c Context) AddPermissions(perms ...Permission) Context {
	return Context{
		perms: append(c.perms, perms...),
	}
}

func (c Context) HasPermission(app string, addr []byte) bool {
	for _, p := range c.perms {
		if app == p.App && bytes.Equal(addr, p.Address) {
			return true
		}
	}
	return false
}

// New should give a fresh context, and know what info makes sense to carry over
func (c Context) New() Context {
	return Context{}
}

// Result captures any non-error wrsp result
// to make sure people use error for error cases
type Result struct {
	Data data.Bytes
	Log  string
}

func (r Result) ToWRSP() wrsp.Result {
	return wrsp.Result{
		Data: r.Data,
		Log:  r.Log,
	}
}
