package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/control"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	"github.com/nosyliam/revolution/pkg/logging"
)

const MainRoutineKind RoutineKind = "Main"

func MainRoutine(macro *Macro) []Action {
	return []Action{
		Log(logging.Info, "Starting Macro"),
	}
}

func init() {
	control.Register(MainRoutineKind, MainRoutine)
}
