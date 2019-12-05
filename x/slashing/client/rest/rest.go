package rest

import (
	"github.com/gorilla/mux"

	"github.com/tepleton/tepleton-sdk/client/context"
	"github.com/tepleton/tepleton-sdk/crypto/keys"
	"github.com/tepleton/tepleton-sdk/wire"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	registerQueryRoutes(ctx, r, cdc)
	registerTxRoutes(ctx, r, cdc, kb)
}
