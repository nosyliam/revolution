package control

import (
	"github.com/nosyliam/revolution/pkg/control/actions"
	"github.com/nosyliam/revolution/pkg/control/common"
)

var Routines = make(map[RoutineKind][]common.Action)

type RoutineKind string

type Routine struct {
	dependencies *common.Dependencies
	actions      []common.Action

	stop  <-chan struct{}
	pause <-chan <-chan struct{}
	err   chan<- string
}

func (r *Routine) Copy(routine *Routine) {
	r.stop = routine.stop
	r.pause = routine.pause
	r.dependencies = routine.dependencies
}

func (r *Routine) Execute() error {
	for _, action := range r.actions {
		if err := action.Execute(r.dependencies); err != nil {
			r.err <- err.Error()
		}
		if len(r.pause) > 0 {
			<-<-r.pause
		}
		if len(r.stop) > 0 {
			<-r.stop
			return nil
		}
	}
}

func NewRoutine(kind RoutineKind, actions ...actions.Action) {
	Routines[kind] = actions
}
