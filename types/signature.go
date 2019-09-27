package types

import crypto "github.com/tepleton/go-crypto"

type StdSignature struct {
	crypto.Signature
	Sequence int64
}
