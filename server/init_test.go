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
	appInit := AppInit{
		AppGenState: mock.AppGenState,
		AppGenTx:    mock.AppGenTx,
	}
	cmd := InitCmd(ctx, cdc, appInit)
	err = cmd.RunE(nil, nil)
	require.NoError(t, err)
}
