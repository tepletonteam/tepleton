package client

import (
	"github.com/tepleton/go-crypto/keys"
	"github.com/tepleton/go-crypto/keys/words"
	dbm "github.com/tepleton/tmlibs/db"
)

// KeyDBName is the directory under root where we store the keys
const KeyDBName = "keys"

// GetKeyManager initializes a key manager based on the configuration
func GetKeyManager(rootDir string) (keys.Keybase, error) {
	db, err := dbm.NewGoLevelDB(KeyDBName, rootDir)
	if err != nil {
		return nil, err
	}
	keybase := keys.New(
		db,
		words.MustLoadCodec("english"),
	)
	return keybase, nil
}
