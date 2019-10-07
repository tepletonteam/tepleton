package mock

import (
	"io/ioutil"
	"os"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/tmlibs/log"
)

// SetupApp returns an application as well as a clean-up function
// to be used to quickly setup a test case with an app
func SetupApp() (wrsp.Application, func(), error) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
		With("module", "mock")
	rootDir, err := ioutil.TempDir("mock-sdk", "")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(rootDir)
	}

	app, err := NewApp(logger, rootDir)
	return app, cleanup, err
}
