package lcd

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tepleton/tmlibs/log"

	tmserver "github.com/tepleton/tepleton/rpc/lib/server"
	cmn "github.com/tepleton/tmlibs/common"

	client "github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/context"
	keys "github.com/tepleton/tepleton-sdk/client/keys"
	rpc "github.com/tepleton/tepleton-sdk/client/rpc"
	tx "github.com/tepleton/tepleton-sdk/client/tx"
	version "github.com/tepleton/tepleton-sdk/version"
	"github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth/client/rest"
	bank "github.com/tepleton/tepleton-sdk/x/bank/client/rest"
	ibc "github.com/tepleton/tepleton-sdk/x/ibc/client/rest"
	stake "github.com/tepleton/tepleton-sdk/x/stake/client/rest"
)

const (
	flagListenAddr = "laddr"
	flagCORS       = "cors"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE:  startRESTServerFn(cdc),
	}
	cmd.Flags().StringP(flagListenAddr, "a", "tcp://localhost:1317", "Address for server to listen on")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	return cmd
}

func startRESTServerFn(cdc *wire.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		listenAddr := viper.GetString(flagListenAddr)
		handler := createHandler(cdc)
		logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
			With("module", "rest-server")
		listener, err := tmserver.StartHTTPServer(listenAddr, handler, logger)
		if err != nil {
			return err
		}

		// Wait forever and cleanup
		cmn.TrapSignal(func() {
			err := listener.Close()
			logger.Error("Error closing listener", "err", err)
		})
		return nil
	}
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", version.RequestHandler).Methods("GET")

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	ctx := context.NewCoreContextFromViper()

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx, r)
	tx.RegisterRoutes(ctx, r, cdc)
	auth.RegisterRoutes(ctx, r, cdc, "acc")
	bank.RegisterRoutes(ctx, r, cdc, kb)
	ibc.RegisterRoutes(ctx, r, cdc, kb)
	stake.RegisterRoutes(ctx, r, cdc, kb)
	return r
}
