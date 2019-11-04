package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/mock"
	"github.com/tepleton/tepleton-sdk/wire"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
)

func TestInit(t *testing.T) {
	defer setupViper(t)()

	logger := log.NewNopLogger()
	cfg, err := tcmd.ParseConfig()
	require.Nil(t, err)
	ctx := NewContext(cfg, logger)
	cdc := wire.NewCodec()
	cmd := InitCmd(ctx, cdc, mock.GenAppParams, nil)
	err = cmd.RunE(nil, nil)
	require.NoError(t, err)
}
