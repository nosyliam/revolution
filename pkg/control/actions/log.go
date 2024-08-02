package actions

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/logging"
	"image"
)

func execArgs(macro *common.Macro, args []interface{}) []interface{} {
	var res []interface{}
	for _, arg := range args {
		switch fn := arg.(type) {
		case func(macro *common.Macro) error:
			res = append(res, fn(macro))
		case func() error:
			res = append(res, fn())
		case common.Action:
			res = append(res, fn.Execute(macro))
		default:
			res = append(res, fn)
		}
	}
	return res
}

type logAction struct {
	log       string
	level     logging.LogLevel
	discord   bool
	verbosity int
	args      []interface{}
	id        *int
}

func (a *logAction) Execute(macro *common.Macro) error {
	msg := fmt.Sprintf(a.log, execArgs(macro, a.args)...)
	if err := macro.Logger.Log(a.verbosity, a.level, msg); err != nil {
		return err
	}
	if a.discord {
		var screenshot *image.RGBA
		if macro.Results.EditedScreenshot != nil {
			screenshot = macro.Results.EditedScreenshot
		} else if macro.Window != nil && macro.Window.Screenshot() != nil {
			screenshot = macro.Window.Screenshot()
		}
		if id, err := macro.Logger.LogDiscord(a.level, msg, a.id, screenshot); err != nil {
			return err
		} else {
			a.id = &id
		}
	}
}

func (a *logAction) Discord(macro *common.Macro) common.Action {
	a.discord = true
	return a
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
