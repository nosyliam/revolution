package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/logging"
)

type logAction struct {
	log       string
	level     logging.LogLevel
	verbosity int
}

func (a *logAction) Execute(deps *common.Macro) error {
	return nil
}

func (a *logAction) V(verbosity int) *logAction {
	a.verbosity = verbosity
	return a
}

func Log(level logging.LogLevel, log string) common.Action {
	return &logAction{level: level, log: log}
}
