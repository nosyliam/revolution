package common

import (
	"context"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
)

var (
	RestartSignal   = errors.New("restart")   // Restarts the routine
	TerminateSignal = errors.New("terminate") // Terminates the routine
	RetrySignal     = errors.New("retry")     // Retry the current action
	StepBackSignal  = errors.New("step back") // Step back to the last action

	AppContext context.Context
)

type RedirectExecution struct {
	Routine RoutineKind
}

func (e RedirectExecution) Error() string { return "redirect" }

type (
	RoutineKind        string
	RoutineExecutor    func(kind RoutineKind)
	SubroutineExecutor func(actions []Action)
)

type Macro struct {
	Account    string
	EventBus   EventBus
	Backend    Backend
	Results    *ActionResults
	State      *config.Object[config.State]
	Settings   *config.Object[config.Settings]
	Database   *config.Object[config.AccountDatabase]
	Logger     *logging.Logger
	Window     *window.Window
	WinManager *window.Manager
	Scratch    *config.Scratch

	Routine    RoutineExecutor
	Subroutine SubroutineExecutor
	Action     func(Action) error
	Status     func(string)
	Pause      <-chan (<-chan struct{})
}

func (m *Macro) Copy() *Macro {
	macro := &Macro{
		EventBus: m.EventBus,
		Backend:  m.Backend,
		Settings: m.Settings,
		Logger:   m.Logger,
		Window:   m.Window,
		State:    m.State,
		Results:  &ActionResults{},
	}
	macro.Action = func(action Action) error {
		return action.Execute(macro)
	}
	return macro
}
