package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/control"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const MainRoutineKind RoutineKind = "Main"

var MainRoutine = Actions{
	Routine(CheckRobloxRoutineKind),
}

func init() {
	control.Register(MainRoutineKind, MainRoutine)
}
