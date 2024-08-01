package control

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
)

var Routines = make(map[common.RoutineKind]common.Actions)

type Routine struct {
	macro   *common.Macro
	actions []common.Action

	depth    int
	redirect *common.RedirectExecution
	parent   *Routine

	stop  <-chan struct{}
	pause <-chan (<-chan struct{})
	err   chan<- string
}

func (r *Routine) Copy(routine *Routine) {
	r.stop = routine.stop
	r.pause = routine.pause
	r.err = routine.err
}

func (r *Routine) Execute() {
	for {
		for i := 0; i < len(r.actions); i++ {
			if err := r.actions[i].Execute(r.macro); err != nil {
				if redirect, ok := err.(*common.RedirectExecution); ok {
					r.parent.redirect = redirect
					return
				}
				switch err {
				case common.RetrySignal:
					i -= 1
					continue
				case common.StepBackSignal:
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
			if r.redirect != nil {
				if r.depth != 0 {
					r.parent.redirect = r.redirect
					return
				}
				r.macro.Routine(r.redirect.Routine)
				break
			}
			if r.macro.State.LoopState.Unwind != nil {
				if r.depth == 0 {
					panic("invalid loop state: not in a loop")
				}
				break
			}
			if len(r.pause) > 0 {
				<-<-r.pause
			}
			if len(r.stop) > 0 {
				<-r.stop
				return
			}
		}
		r.macro.Results = &common.ActionResults{}
		if r.depth != 0 {
			break
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
	var execSub func(routine *Routine, macro *common.Macro) common.SubroutineExecutor
	var exec func(routine *Routine, macro *common.Macro) common.RoutineExecutor
	execSub = func(routine *Routine, macro *common.Macro) common.SubroutineExecutor {
		return func(actions []common.Action) {
			subMacro := macro.Copy()
			subMacro.Logger = macro.Logger
			subMacro.Status = macro.Status
			subMacro.Routine = exec(routine, subMacro)
			subMacro.Subroutine = execSub(routine, subMacro)
			subMacro.Action = func(action common.Action) error {
				return action.Execute(subMacro)
			}
			subRoutine := &Routine{macro: subMacro, actions: actions, depth: routine.depth + 1, parent: routine.parent}
			subRoutine.Copy(routine)
			subRoutine.Execute()
		}
	}
	exec = func(routine *Routine, macro *common.Macro) common.RoutineExecutor {
		return func(kind common.RoutineKind) {
			subActions, ok := Routines[kind]
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
			subRoutine := &Routine{macro: subMacro, actions: subActions, depth: routine.depth + 1, parent: routine}
			subRoutine.Copy(routine)
			subRoutine.Execute()
		}
	}
	macro.Routine = exec(routine, macro)
	macro.Subroutine = execSub(routine, macro)
	macro.Action = func(action common.Action) error {
		return action.Execute(macro)
	}
	macro.Status = func(stat string) {
		status <- stat
	}
	routine.Execute()
}

func Register(kind common.RoutineKind, fn common.Actions) {
	Routines[kind] = fn
}
