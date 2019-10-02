package app

import (
	bam "github.com/tepleton/tepleton-sdk/baseapp"
)

type testBasecoinApp struct {
	*BasecoinApp
	*bam.TestApp
}

func newTestBasecoinApp() *testBasecoinApp {
	app := NewBasecoinApp("")
	tba := &testBasecoinApp{
		BasecoinApp: app,
	}
	tba.TestApp = bam.NewTestApp(app.BaseApp)
	return tba
}
