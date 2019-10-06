package main

import (
	"fmt"
	"os"

	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
)

func main() {
	fmt.Println("This is temporary, for unblocking our build process.")
	return

	// TODO CREATE CLI
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	db, err := dbm.NewGoLevelDB("basecoind", "data")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bapp := app.NewBasecoinApp(logger, db)
	baseapp.RunForever(bapp)
}
