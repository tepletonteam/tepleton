package main

import "github.com/tepleton/tepleton-sdk/examples/basecoin/app"

func main() {
	// TODO CREATE CLI

	bapp := app.NewBasecoinApp("")
	bapp.RunForever()
}
