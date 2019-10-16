package server

import (
	"github.com/tepleton/go-crypto/keys"
	"github.com/tepleton/go-crypto/keys/words"
	dbm "github.com/tepleton/tmlibs/db"

	sdk "github.com/tepleton/tepleton-sdk/types"
)

// GenerateCoinKey returns the address of a public key,
// along with the secret phrase to recover the private key.
// You can give coins to this address and return the recovery
// phrase to the user to access them.
func GenerateCoinKey() (sdk.Address, string, error) {
	// construct an in-memory key store
	codec, err := words.LoadCodec("english")
	if err != nil {
		return nil, "", err
	}
	keybase := keys.New(
		dbm.NewMemDB(),
		codec,
	)

	// generate a private key, with recovery phrase
	info, secret, err := keybase.Create("name", "pass", keys.AlgoEd25519)
	if err != nil {
		return nil, "", err
	}

	addr := info.PubKey.Address()
	return addr, secret, nil
}
