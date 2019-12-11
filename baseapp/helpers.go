package baseapp

import (
	"github.com/tepleton/tepleton/wrsp/server"
	wrsp "github.com/tepleton/tepleton/wrsp/types"
	cmn "github.com/tepleton/tepleton/libs/common"
)

// RunForever - BasecoinApp execution and cleanup
func RunForever(app wrsp.Application) {

	// Start the WRSP server
	srv, err := server.NewServer("0.0.0.0:26658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
		return
	}
	err = srv.Start()
	if err != nil {
		cmn.Exit(err.Error())
		return
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		err := srv.Stop()
		if err != nil {
			cmn.Exit(err.Error())
		}
	})
}
