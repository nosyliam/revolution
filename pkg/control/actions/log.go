package actions

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/pkg/errors"
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

type LogAction struct {
	log       string
	level     logging.LogLevel
	discord   bool
	status    bool
	verbosity int
	args      []interface{}
	id        *int
}

func (a *LogAction) Execute(macro *common.Macro) error {
	if a.status {
		macro.Status(a.log)
	}
	msg := fmt.Sprintf(a.log, execArgs(macro, a.args)...)
	if err := macro.Logger.Log(a.verbosity, a.level, msg); err != nil {
		return errors.Wrap(err, "failed to log")
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
	return nil
}

func (a *LogAction) Status() *LogAction {
	a.status = true
	return a
}

func (a *LogAction) Discord() *LogAction {
	a.discord = true
	return a
}

func (a *LogAction) V(verbosity int) *LogAction {
	a.verbosity = verbosity
	return &LogAction{verbosity: verbosity, level: a.level, log: a.log}
}

func Error(log string, args ...interface{}) *LogAction {
	return &LogAction{level: logging.Error, log: log, args: args}
}

func Info(log string, args ...interface{}) *LogAction {
	return &LogAction{level: logging.Info, log: log, args: args}
}

func Warning(log string, args ...interface{}) *LogAction {
	return &LogAction{level: logging.Warning, log: log, args: args}
}
