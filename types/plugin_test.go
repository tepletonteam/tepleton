package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
)

//----------------------------------

type Dummy struct{}

func (d *Dummy) Name() string {
	return "dummy"
}
func (d *Dummy) RunTx(store KVStore, ctx CallContext, txBytes []byte) (res wrsp.Result) {
	return
}
func (d *Dummy) SetOption(storei KVStore, key, value string) (log string) {
	return ""
}
func (d *Dummy) InitChain(store KVStore, vals []*wrsp.Validator) {
}
func (d *Dummy) BeginBlock(store KVStore, hash []byte, header *wrsp.Header) {
}
func (d *Dummy) EndBlock(store KVStore, height uint64) (res wrsp.ResponseEndBlock) {
	return
}

//----------------------------------

func TestPlugin(t *testing.T) {
	assert := assert.New(t)
	plugins := NewPlugins()
	assert.Zero(len(plugins.GetList()), "plugins object init with a objects")
	plugins.RegisterPlugin(&Dummy{})
	assert.Equal(len(plugins.GetList()), 1, "plugin wasn't added to plist after registered")
	assert.Equal(plugins.GetByName("dummy").Name(), "dummy", "plugin wasn't retrieved properly with GetByName")
}
