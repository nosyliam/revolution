package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type routineAction struct {
	name common.RoutineKind
}

func (a *routineAction) Execute(macro *common.Macro) error {
	macro.Routine(a.name)
	return nil
}

func Routine(kind common.RoutineKind) common.Action {
	return &routineAction{name: kind}
}

type subroutineAction struct {
	actions []common.Action
}

func (a *subroutineAction) Execute(macro *common.Macro) error {
	macro.Subroutine(a.actions)
	return nil
}

func Subroutine(actions ...common.Action) common.Action {
	return &subroutineAction{actions: actions}
}

type statusAction struct {
	status string
}

func (a *statusAction) Execute(deps *common.Macro) error {
	deps.Status(a.status)
	return nil
}

func Status(status string) common.Action {
	return &statusAction{status: status}
}
