package server

import (
	"encoding/json"
	"path/filepath"

	wrsp "github.com/tepleton/tepleton/wrsp/types"
	tmtypes "github.com/tepleton/tepleton/types"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"
)

// AppCreator lets us lazily initialize app, using home dir
// and other flags (?) to start
type AppCreator func(string, log.Logger) (wrsp.Application, error)

// AppExporter dumps all app state to JSON-serializable structure and returns the current validator set
type AppExporter func(home string, log log.Logger) (json.RawMessage, []tmtypes.GenesisValidator, error)

// ConstructAppCreator returns an application generation function
func ConstructAppCreator(appFn func(log.Logger, dbm.DB) wrsp.Application, name string) AppCreator {
	return func(rootDir string, logger log.Logger) (wrsp.Application, error) {
		dataDir := filepath.Join(rootDir, "data")
		db, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, err
		}
		app := appFn(logger, db)
		return app, nil
	}
}

// ConstructAppExporter returns an application export function
func ConstructAppExporter(appFn func(log.Logger, dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error), name string) AppExporter {
	return func(rootDir string, logger log.Logger) (json.RawMessage, []tmtypes.GenesisValidator, error) {
		dataDir := filepath.Join(rootDir, "data")
		db, err := dbm.NewGoLevelDB(name, dataDir)
		if err != nil {
			return nil, nil, err
		}
		return appFn(logger, db)
	}
}
