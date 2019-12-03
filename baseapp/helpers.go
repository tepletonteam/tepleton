package baseapp

import (
	"github.com/tepleton/wrsp/server"
	wrsp "github.com/tepleton/wrsp/types"
	cmn "github.com/tepleton/tmlibs/common"
)

// RunForever - BasecoinApp execution and cleanup
func RunForever(app wrsp.Application) {

	// Start the WRSP server
	srv, err := server.NewServer("0.0.0.0:26658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
	}
	srv.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}
