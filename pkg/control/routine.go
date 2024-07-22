package control

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
)

var Routines = make(map[common.RoutineKind]common.RoutineFunc)

const MainRoutine common.RoutineKind = "main"

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
	for _, action := range r.actions {
		if err := action.Execute(r.macro); err != nil {
			r.err <- err.Error()
			<-<-r.pause
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
			subMacro.SetStatus = macro.SetStatus
			subMacro.ExecRoutine = exec(routine, subMacro)
			subMacro.ExecAction = func(action common.Action) error {
				return action.Execute(subMacro)
			}
			subRoutine := &Routine{macro: subMacro, actions: fn(subMacro)}
			subRoutine.Copy(routine)
			subRoutine.Execute()
		}
	}
	macro.ExecRoutine = exec(routine, macro)
	macro.ExecAction = func(action common.Action) error {
		return action.Execute(macro)
	}
	macro.SetStatus = func(stat string) {
		status <- stat
	}
	routine.Execute()
}

func Register(kind common.RoutineKind, fn common.RoutineFunc) {
	Routines[kind] = fn
}
