package crypto

import (
	amino "github.com/tepleton/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	// NOTE: It's important that there be no conflicts here,
	// as that would change the canonical representations,
	// and therefore change the address.
	// TODO: Add feature to go-amino to ensure that there
	// are no conflicts.
	RegisterAmino(cdc)
}

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*PubKey)(nil), nil)
	cdc.RegisterConcrete(PubKeyEd25519{},
		"tepleton/PubKeyEd25519", nil)
	cdc.RegisterConcrete(PubKeySecp256k1{},
		"tepleton/PubKeySecp256k1", nil)

	cdc.RegisterInterface((*PrivKey)(nil), nil)
	cdc.RegisterConcrete(PrivKeyEd25519{},
		"tepleton/PrivKeyEd25519", nil)
	cdc.RegisterConcrete(PrivKeySecp256k1{},
		"tepleton/PrivKeySecp256k1", nil)

	cdc.RegisterInterface((*Signature)(nil), nil)
	cdc.RegisterConcrete(SignatureEd25519{},
		"tepleton/SignatureKeyEd25519", nil)
	cdc.RegisterConcrete(SignatureSecp256k1{},
		"tepleton/SignatureKeySecp256k1", nil)
}
