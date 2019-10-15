package client_test

import (
	"os"
	"testing"

	nm "github.com/tepleton/tepleton/node"
)

var node *nm.Node

// See https://golang.org/pkg/testing/#hdr-Main
// for more details
func TestMain(m *testing.M) {
	// start a basecoind node and LCD server in the background to test against

	// run all the tests against a single server instance
	code := m.Run()

	// tear down

	//
	os.Exit(code)
}
