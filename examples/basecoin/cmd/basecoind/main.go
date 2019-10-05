package main

import (
	"github.com/tepleton/tepleton-sdk/baseapp"
	"github.com/tepleton/tepleton-sdk/examples/basecoin/app"
)

func main() {
	// TODO CREATE CLI

	bapp := app.NewBasecoinApp("")
	baseapp.RunForever(bapp)
}
