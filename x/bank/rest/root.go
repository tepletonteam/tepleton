package rest

import (
	"github.com/gorilla/mux"

	keys "github.com/tepleton/go-crypto/keys"

	"github.com/tepleton/tepleton-sdk/wire"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", SendRequestHandler(cdc, kb)).Methods("POST")
}
