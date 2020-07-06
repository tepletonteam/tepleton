package client

import (
	"github.com/tepleton/go-crypto/keys"
	"github.com/tepleton/go-crypto/keys/words"
	dbm "github.com/tepleton/tmlibs/db"
)

// GetKeyBase initializes a keybase based on the given db.
// The KeyBase manages all activity requiring access to a key.
func GetKeyBase(db dbm.DB) keys.Keybase {
	keybase := keys.New(
		db,
		words.MustLoadCodec("english"),
	)
	return keybase
}

// MockKeyBase generates an in-memory keybase that will be discarded
// useful for --dry-run to generate a seed phrase without
// storing the key
func MockKeyBase() keys.Keybase {
	return GetKeyBase(dbm.NewMemDB())
}
