package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"

	wrsp "github.com/tepleton/wrsp/types"
	wire "github.com/tepleton/go-wire"

	"github.com/tepleton/basecoin/types"

	client "github.com/tepleton/tepleton/rpc/client"
	tmtypes "github.com/tepleton/tepleton/types"
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

// StripHex remove the first two hex bytes
func StripHex(s string) string {
	if isHex(s) {
		return s[2:]
	}
	return s
}

// Query - send an wrsp query
func Query(tmAddr string, key []byte) (*wrsp.ResultQuery, error) {
	httpClient := client.NewHTTP(tmAddr, "/websocket")
	return queryWithClient(httpClient, key)
}

func queryWithClient(httpClient *client.HTTP, key []byte) (*wrsp.ResultQuery, error) {
	res, err := httpClient.WRSPQuery("/key", key, true)
	if err != nil {
		return nil, errors.Errorf("Error calling /wrsp_query: %v", err)
	}
	if !res.Code.IsOK() {
		return nil, errors.Errorf("Query got non-zero exit code: %v. %s", res.Code, res.Log)
	}
	return res.ResultQuery, nil
}

// fetch the account by querying the app
func getAccWithClient(httpClient *client.HTTP, address []byte) (*types.Account, error) {

	key := types.AccountKey(address)
	response, err := queryWithClient(httpClient, key)
	if err != nil {
		return nil, err
	}

	accountBytes := response.Value

	if len(accountBytes) == 0 {
		return nil, fmt.Errorf("Account bytes are empty for address: %X ", address) //never stack trace
	}

	var acc *types.Account
	err = wire.ReadBinaryBytes(accountBytes, &acc)
	if err != nil {
		return nil, errors.Errorf("Error reading account %X error: %v",
			accountBytes, err.Error())
	}

	return acc, nil
}

func getHeaderAndCommit(tmAddr string, height int) (*tmtypes.Header, *tmtypes.Commit, error) {
	httpClient := client.NewHTTP(tmAddr, "/websocket")
	res, err := httpClient.Commit(height)
	if err != nil {
		return nil, nil, errors.Errorf("Error on commit: %v", err)
	}
	header := res.Header
	commit := res.Commit

	return header, commit, nil
}

func waitForBlock(httpClient *client.HTTP) error {
	res, err := httpClient.Status()
	if err != nil {
		return err
	}

	lastHeight := res.LatestBlockHeight
	for {
		res, err := httpClient.Status()
		if err != nil {
			return err
		}
		if res.LatestBlockHeight > lastHeight {
			break
		}

	}
	return nil
}
