package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/urfave/cli"

	"github.com/tepleton/basecoin/types"
	cmn "github.com/tepleton/go-common"
	client "github.com/tepleton/go-rpc/client"
	"github.com/tepleton/go-wire"
	ctypes "github.com/tepleton/tepleton/rpc/core/types"
)

func cmdAccount(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("account command requires an argument ([address])")
	}
	addrHex := c.Args()[0]

	// convert destination address to bytes
	addr, err := hex.DecodeString(addrHex)
	if err != nil {
		return errors.New("Account address is invalid hex: " + err.Error())
	}

	acc, err := getAcc(c, addr)
	if err != nil {
		return err
	}
	fmt.Println(string(wire.JSONBytes(acc)))
	return nil
}

// fetch the account by querying the app
func getAcc(c *cli.Context, address []byte) (*types.Account, error) {
	tmAddr := c.String("tepleton")
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
