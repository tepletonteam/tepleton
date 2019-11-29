package lcd

import (
	"fmt"
	"net/http"

	"github.com/tepleton/tepleton-sdk/client/context"
	"github.com/tepleton/tepleton-sdk/version"
	"github.com/tepleton/tepleton-sdk/wire"
)

// cli version REST handler endpoint
func CLIVersionRequestHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()
	w.Write([]byte(v))
}

// connected node version REST handler endpoint
func NodeVersionRequestHandler(cdc *wire.Codec, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := ctx.Query("/app/version")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query version. Error: %s", err.Error())))
			return
		}
		w.Write([]byte(version))
	}
}
