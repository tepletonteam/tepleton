package keys

import (
	amino "github.com/tepleton/go-amino"
	crypto "github.com/tepleton/go-crypto"
)

var cdc = amino.NewCodec()

func init() {
	crypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
	cdc.RegisterConcrete(ledgerInfo{}, "crypto/keys/ledgerInfo", nil)
	cdc.RegisterConcrete(offlineInfo{}, "crypto/keys/offlineInfo", nil)
}
