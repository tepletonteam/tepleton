package auth

import (
	crypto "github.com/tepleton/go-crypto"
)

type SetPubKeyer interface {
	SetPubKey(crypto.PubKey)
}
