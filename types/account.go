package types

import (
	"encoding/hex"
	"errors"
	"fmt"

	bech32tepleton "github.com/tepleton/bech32tepleton/go"
	crypto "github.com/tepleton/go-crypto"
	cmn "github.com/tepleton/tmlibs/common"
)

//Address is a go crypto-style Address
type Address = cmn.HexBytes

// Bech32 prefixes
const (
	Bech32PrefixAccAddr = "tepletonaccaddr"
	Bech32PrefixAccPub  = "tepletonaccpub"
	Bech32PrefixValAddr = "tepletonvaladdr"
	Bech32PrefixValPub  = "tepletonvalpub"
)

// Bech32TepletonifyAcc takes Address and returns the Bech32Tepleton encoded string
func Bech32TepletonifyAcc(addr Address) (string, error) {
	return bech32tepleton.ConvertAndEncode(Bech32PrefixAccAddr, addr.Bytes())
}

// Bech32TepletonifyAccPub takes AccountPubKey and returns the Bech32Tepleton encoded string
func Bech32TepletonifyAccPub(pub crypto.PubKey) (string, error) {
	return bech32tepleton.ConvertAndEncode(Bech32PrefixAccPub, pub.Bytes())
}

// Bech32TepletonifyVal returns the Bech32Tepleton encoded string for a validator address
func Bech32TepletonifyVal(addr Address) (string, error) {
	return bech32tepleton.ConvertAndEncode(Bech32PrefixValAddr, addr.Bytes())
}

// Bech32TepletonifyValPub returns the Bech32Tepleton encoded string for a validator pubkey
func Bech32TepletonifyValPub(pub crypto.PubKey) (string, error) {
	return bech32tepleton.ConvertAndEncode(Bech32PrefixValPub, pub.Bytes())
}

// create an Address from a string
func GetAccAddressHex(address string) (addr Address, err error) {
	if len(address) == 0 {
		return addr, errors.New("must use provide address")
	}
	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}

// create an Address from a string
func GetAccAddressBech32Tepleton(address string) (addr Address, err error) {
	bz, err := getFromBech32Tepleton(address, Bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}

// create an Address from a hex string
func GetValAddressHex(address string) (addr Address, err error) {
	if len(address) == 0 {
		return addr, errors.New("must use provide address")
	}
	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}

// create an Address from a bech32tepleton string
func GetValAddressBech32Tepleton(address string) (addr Address, err error) {
	bz, err := getFromBech32Tepleton(address, Bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}

//Decode a validator publickey into a public key
func GetValPubKeyBech32Tepleton(pubkey string) (pk crypto.PubKey, err error) {
	bz, err := getFromBech32Tepleton(pubkey, Bech32PrefixValPub)
	if err != nil {
		return nil, err
	}

	pk, err = crypto.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func getFromBech32Tepleton(bech32, prefix string) ([]byte, error) {
	if len(bech32) == 0 {
		return nil, errors.New("must provide non-empty string")
	}
	hrp, bz, err := bech32tepleton.DecodeAndConvert(bech32)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("Invalid bech32 prefix. Expected %s, Got %s", prefix, hrp)
	}

	return bz, nil
}
