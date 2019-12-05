package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/slashing"
)

func registerQueryRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc(
		"/slashing/signing_info/{validator}",
		signingInfoHandlerFn(ctx, "slashing", cdc),
	).Methods("GET")
}

// http request handler to query signing info
func signingInfoHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read parameters
		vars := mux.Vars(r)
		bech32validator := vars["validator"]

		validatorAddr, err := sdk.GetValAddressBech32(bech32validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key := slashing.GetValidatorSigningInfoKey(validatorAddr)
		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query signing info. Error: %s", err.Error())))
			return
		}

		var signingInfo slashing.ValidatorSigningInfo
		err = cdc.UnmarshalBinary(res, &signingInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't decode signing info. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(signingInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
