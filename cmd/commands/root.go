package commands

import (
	"os"

	"github.com/tepleton/tmlibs/log"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)
