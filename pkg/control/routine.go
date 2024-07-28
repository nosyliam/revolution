package control

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
)

var Routines = make(map[common.RoutineKind]common.RoutineFunc)

type Routine struct {
	macro   *common.Macro
	actions []common.Action

	stop  <-chan struct{}
	pause <-chan (<-chan struct{})
	err   chan<- string
}

func (r *Routine) Copy(routine *Routine) {
	r.stop = routine.stop
	r.pause = routine.pause
	r.macro = routine.macro
}

func (r *Routine) Execute() {
	for i := 0; i < len(r.actions); i++ {
		if err := r.actions[i].Execute(r.macro); err != nil {
			switch err {
			case common.RetrySignal:
				i -= 2
				continue
			case common.RestartSignal:
				i = 0
				continue
			case common.TerminateSignal:
				return
			default:
				r.err <- err.Error()
				<-<-r.pause
			}
		}
		if len(r.pause) > 0 {
			<-<-r.pause
		}
		if len(r.stop) > 0 {
			<-r.stop
		}
	}
}

func ExecuteRoutine(
	macro *common.Macro,
	actions []common.Action,
	stop <-chan struct{},
	pause <-chan (<-chan struct{}),
	status chan<- string,
	err chan<- string,
) {
	routine := &Routine{
		macro:   macro,
		actions: actions,
		stop:    stop,
		pause:   pause,
		err:     err,
	}
	var exec func(routine *Routine, macro *common.Macro) common.RoutineExecutor
	exec = func(routine *Routine, macro *common.Macro) common.RoutineExecutor {
		return func(kind common.RoutineKind) {
			fn, ok := Routines[kind]
			if !ok {
				panic(fmt.Sprintf("unknown subroutine %s", string(kind)))
			}
			subMacro := macro.Copy()
			subMacro.Logger = subMacro.Logger.Child(string(kind))
			subMacro.Results = &common.ActionResults{}
			subMacro.Status = macro.Status
			subMacro.Routine = exec(routine, subMacro)
			subMacro.Action = func(action common.Action) error {
				return action.Execute(subMacro)
			}
			subRoutine := &Routine{macro: subMacro, actions: fn(subMacro)}
			subRoutine.Copy(routine)
			subRoutine.Execute()
		}
	}
	macro.Routine = exec(routine, macro)
	macro.Action = func(action common.Action) error {
		return action.Execute(macro)
	}
	macro.Status = func(stat string) {
		status <- stat
	}
	routine.Execute()
}

func Register(kind common.RoutineKind, fn common.RoutineFunc) {
	Routines[kind] = fn
}
