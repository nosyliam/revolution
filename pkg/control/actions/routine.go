package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type routineAction struct {
	name common.RoutineKind
}

func (a *routineAction) Execute(deps *common.Macro) error {
	deps.ExecRoutine(a.name)
	return nil
}

func Routine(kind common.RoutineKind) common.Action {
	return &routineAction{name: kind}
}

type statusAction struct {
	status string
}

func (a *statusAction) Execute(deps *common.Macro) error {
	deps.SetStatus(a.status)
	return nil
}

func Status(status string) common.Action {
	return &statusAction{status: status}
}
