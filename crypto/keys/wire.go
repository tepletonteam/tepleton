package keys

import (
	ccrypto "github.com/tepleton/tepleton-sdk/crypto"
	amino "github.com/tepleton/go-amino"
	tcrypto "github.com/tepleton/tepleton/crypto"
)

var cdc = amino.NewCodec()

func init() {
	tcrypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(ccrypto.PrivKeyLedgerSecp256k1{},
		"tepleton/PrivKeyLedgerSecp256k1", nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
	cdc.RegisterConcrete(ledgerInfo{}, "crypto/keys/ledgerInfo", nil)
	cdc.RegisterConcrete(offlineInfo{}, "crypto/keys/offlineInfo", nil)
}
