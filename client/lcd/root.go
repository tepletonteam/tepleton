package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	client "github.com/tepleton/tepleton-sdk/client"
	keys "github.com/tepleton/tepleton-sdk/client/keys"
	rpc "github.com/tepleton/tepleton-sdk/client/rpc"
	tx "github.com/tepleton/tepleton-sdk/client/tx"
	version "github.com/tepleton/tepleton-sdk/version"
	"github.com/tepleton/tepleton-sdk/wire"
)

const (
	flagBind = "bind"
	flagCORS = "cors"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE:  startRESTServer(cdc),
	}
	// TODO: handle unix sockets also?
	cmd.Flags().StringP(flagBind, "b", "localhost:1317", "Interface and port that server binds to")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	return cmd
}

func startRESTServer(cdc *wire.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		bind := viper.GetString(flagBind)
		r := initRouter(cdc)
		return http.ListenAndServe(bind, r)
	}
}

func initRouter(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", version.VersionRequestHandler).Methods("GET")

	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(r)
	tx.RegisterRoutes(r, cdc)
	return r
}
