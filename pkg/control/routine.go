package control

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/pkg/errors"
)

type Routine struct {
	macro   *common.Macro
	actions []common.Action

	depth       int
	redirectLoc *common.RedirectExecution
	parent      *Routine

	stop     chan struct{}
	redirect chan *common.RedirectExecution
	err      chan<- string
}

func (r *Routine) Copy(routine *Routine) {
	r.stop = routine.stop
	r.err = routine.err
	r.redirect = routine.redirect
}

func (r *Routine) Execute() {
	for {
		for i := 0; i < len(r.actions); i++ {
			if err := r.actions[i].Execute(r.macro); err != nil {
				if redirect, ok := err.(*common.RedirectExecution); ok {
					r.parent.redirectLoc = redirect
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
					r.err <- errors.Wrap(err, "").Error()
					<-<-r.macro.Pause
				}
			}
			if len(r.redirect) > 0 {
				r.parent.redirectLoc = <-r.redirect
				return
			}
			if len(r.stop) > 0 {
				<-r.stop
				if r.depth != 0 {
					r.stop <- struct{}{}
				}
				return
			}
			if r.redirectLoc != nil {
				if r.depth != 0 {
					r.parent.redirectLoc = r.redirectLoc
					return
				}
				r.macro.Routine(r.redirectLoc.Routine)
				break
			}
			if r.macro.Scratch.LoopState.Unwind != nil {
				if r.depth == 0 {
					panic("invalid loop state: not in a loop")
				}
				break
			}
			if len(r.macro.Pause) > 0 {
				<-<-r.macro.Pause
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
	actions common.Actions,
	stop chan struct{},
	status chan<- string,
	err chan<- string,
	redirect chan *common.RedirectExecution,
) {
	routine := &Routine{
		macro:    macro,
		actions:  actions,
		stop:     stop,
		err:      err,
		redirect: redirect,
	}
	macro.Scratch.Stack = []string{"Main"}
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
			subActions, ok := common.Routines[kind]
			if !ok {
				panic(fmt.Sprintf("unknown subroutine %s", string(kind)))
			}
			macro.Scratch.Stack = append([]string{string(kind)}, macro.Scratch.Stack...)
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
			macro.Scratch.Stack = macro.Scratch.Stack[1:]
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
