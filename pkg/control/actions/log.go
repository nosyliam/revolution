package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/logging"
)

func execArgs(args []interface{}) []interface{} {
	var res []interface{}
}

type logAction struct {
	log       string
	level     logging.LogLevel
	discord   bool
	verbosity int
	args      []interface{}
}

func (a *logAction) Execute(deps *common.Macro) error {
	return nil
}

func (a *logAction) V(verbosity int) *logAction {
	a.verbosity = verbosity
	return &logAction{verbosity: verbosity, level: a.level, log: a.log}
}

func Error(log string, args ...interface{}) common.Action {
	return &logAction{level: logging.Error, log: log, args: args}
}

type statusAction struct {
	msg  string
	args []interface{}
}

func (a *statusAction) Execute(deps *common.Macro) error {
	return nil
}

func Status(status string, args ...interface{}) common.Action {
	return &statusAction{msg: status, args: args}
}
