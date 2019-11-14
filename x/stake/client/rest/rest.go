package rest

import (
	"github.com/gorilla/mux"
	"github.com/tepleton/go-crypto/keys"

	"github.com/tepleton/tepleton-sdk/client/context"
	"github.com/tepleton/tepleton-sdk/wire"
)

func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	registerQueryRoutes(ctx, r, cdc, kb)
	registerTxRoutes(ctx, r, cdc, kb)
}
