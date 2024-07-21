package common

import (
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/nosyliam/revolution/pkg/window"
)

type (
	RoutineKind     string
	RoutineFunc     func(macro *Macro) []Action
	RoutineExecutor func(kind RoutineKind) error
)

type Macro struct {
	EventBus    EventBus
	Backend     Backend
	Results     *ActionResults
	Settings    *config.Settings
	Logger      *logging.Logger
	Window      *window.Window
	State       *config.MacroState
	ExecRoutine RoutineExecutor
	ExecAction  func(Action) error
}

func (m *Macro) Copy() *Macro {
	macro := &Macro{
		EventBus: m.EventBus,
		Backend:  m.Backend,
		Settings: m.Settings,
		Logger:   m.Logger,
		Window:   m.Window,
		State:    m.State,
	}
	macro.ExecAction = func(action Action) error {
		return action.Execute(macro)
	}
	return macro
}
