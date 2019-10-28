package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/tepleton-sdk/mock"
)

func TestInit(t *testing.T) {
	defer setupViper(t)()

	logger := log.NewNopLogger()
	cmd := InitCmd(mock.GenInitOptions, logger)
	err := cmd.RunE(nil, nil)
	require.NoError(t, err)
}
