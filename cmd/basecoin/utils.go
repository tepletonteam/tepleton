package main

import (
	"encoding/hex"
	"errors"

	"github.com/tepleton/basecoin/types"

	cmn "github.com/tepleton/go-common"
	client "github.com/tepleton/go-rpc/client"
	"github.com/tepleton/go-wire"
	ctypes "github.com/tepleton/tepleton/rpc/core/types"
)

// Returns true for non-empty hex-string prefixed with "0x"
func isHex(s string) bool {
	if len(s) > 2 && s[:2] == "0x" {
		_, err := hex.DecodeString(s[2:])
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func stripHex(s string) string {
	if isHex(s) {
		return s[2:]
	}
	return s
}

// fetch the account by querying the app
func getAcc(tmAddr string, address []byte) (*types.Account, error) {
	clientURI := client.NewClientURI(tmAddr)
	tmResult := new(ctypes.TMResult)

	params := map[string]interface{}{
		"path":  "/key",
		"data":  append([]byte("base/a/"), address...),
		"prove": false,
	}
	_, err := clientURI.Call("wrsp_query", params, tmResult)
	if err != nil {
		return nil, errors.New(cmn.Fmt("Error calling /wrsp_query: %v", err))
	}
	res := (*tmResult).(*ctypes.ResultWRSPQuery)
	if !res.Response.Code.IsOK() {
		return nil, errors.New(cmn.Fmt("Query got non-zero exit code: %v. %s", res.Response.Code, res.Response.Log))
	}
	accountBytes := res.Response.Value

	if len(accountBytes) == 0 {
		return nil, errors.New(cmn.Fmt("Account bytes are empty from query for address %X", address))
	}
	var acc *types.Account
	err = wire.ReadBinaryBytes(accountBytes, &acc)
	if err != nil {
		return nil, errors.New(cmn.Fmt("Error reading account %X error: %v",
			accountBytes, err.Error()))
	}

	return acc, nil
}
