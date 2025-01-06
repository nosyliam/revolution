package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/movement"
)

type SleepAction struct {
	ms      int
	seconds bool
}

func (a *SleepAction) Seconds() *SleepAction {
	a.seconds = true
	return a
}

func (a *SleepAction) Execute(macro *common.Macro) error {
	var ms = a.ms
	if a.seconds {
		ms *= 1000
	}
	movement.Sleep(ms, macro)
	return nil
}

func Sleep(ms int) *SleepAction {
	return &SleepAction{ms, false}
}

type executePatternAction struct {
	name interface{}
	id   string
}

func (e *executePatternAction) Execute(macro *common.Macro) error {
	var computed string
	switch val := e.name.(type) {
	case func(macro *common.Macro) string:
		computed = val(macro)
	case string:
		computed = val
	}
	return macro.Pattern.Execute(macro, nil, computed)
}

func ExecutePattern(name interface{}) common.Action {
	return &executePatternAction{name: name}
}
