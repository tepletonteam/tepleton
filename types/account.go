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

// Bech32TepletonifyAcc takes Address and returns the Bech32Tepleton encoded string
func Bech32TepletonifyAcc(addr Address) (string, error) {
	return bech32tepleton.ConvertAndEncode("tepletonaccaddr", addr.Bytes())
}

// Bech32TepletonifyAccPub takes AccountPubKey and returns the Bech32Tepleton encoded string
func Bech32TepletonifyAccPub(pub crypto.PubKey) (string, error) {
	return bech32tepleton.ConvertAndEncode("tepletonaccpub", pub.Bytes())
}

// Bech32TepletonifyVal returns the Bech32Tepleton encoded string for a validator address
func Bech32TepletonifyVal(addr Address) (string, error) {
	return bech32tepleton.ConvertAndEncode("tepletonvaladdr", addr.Bytes())
}

// Bech32TepletonifyValPub returns the Bech32Tepleton encoded string for a validator pubkey
func Bech32TepletonifyValPub(pub crypto.PubKey) (string, error) {
	return bech32tepleton.ConvertAndEncode("tepletonvalpub", pub.Bytes())
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
	if len(address) == 0 {
		return addr, errors.New("must use provide address")
	}

	hrp, bz, err := bech32tepleton.DecodeAndConvert(address)

	if hrp != "tepletonaccaddr" {
		return addr, fmt.Errorf("Invalid Address Prefix. Expected tepletonaccaddr, Got %s", hrp)
	}

	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}

// create an Address from a string
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

// create an Address from a string
func GetValAddressBech32Tepleton(address string) (addr Address, err error) {
	if len(address) == 0 {
		return addr, errors.New("must use provide address")
	}

	hrp, bz, err := bech32tepleton.DecodeAndConvert(address)

	if hrp != "tepletonvaladdr" {
		return addr, fmt.Errorf("Invalid Address Prefix. Expected tepletonvaladdr, Got %s", hrp)
	}

	if err != nil {
		return nil, err
	}
	return Address(bz), nil
}
