package rest

import (
	"github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth/commands"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc("/accounts/{address}", QueryAccountRequestHandler(storeName, cdc, auth.GetAccountDecoder(cdc))).Methods("GET")
}
