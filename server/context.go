package server

import (
	cfg "github.com/tepleton/tepleton/config"
	"github.com/tepleton/tmlibs/log"
)

type Context struct {
	Config *cfg.Config
	Logger log.Logger
}

func NewContext(config *cfg.Config, logger log.Logger) *Context {
	return &Context{config, logger}
}
