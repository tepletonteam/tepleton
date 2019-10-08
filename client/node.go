package client

import rpcclient "github.com/tepleton/tepleton/rpc/client"

// GetNode prepares a simple rpc.Client from the flags
func GetNode(uri string) rpcclient.Client {
	return rpcclient.NewHTTP(uri, "/websocket")
}
