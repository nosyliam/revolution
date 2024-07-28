package common

import (
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
)

var (
	RestartSignal   = errors.New("restart")   // RestartSignal the routine
	RetrySignal     = errors.New("retry")     // Step back and retry the last action
	TerminateSignal = errors.New("terminate") // TerminateSignal execution of the routine
)

type (
	RoutineKind     string
	RoutineFunc     func(macro *Macro) []Action
	RoutineExecutor func(kind RoutineKind)
)

type Macro struct {
	EventBus   EventBus
	Backend    Backend
	Results    *ActionResults
	Settings   *config.Settings
	Logger     *logging.Logger
	Window     *window.Window
	WinManager *window.Manager
	State      *config.MacroState

	Routine RoutineExecutor
	Action  func(Action) error
	Status  func(string)
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
	macro.Action = func(action Action) error {
		return action.Execute(macro)
	}
	return macro
}
