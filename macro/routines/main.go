package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const MainRoutineKind RoutineKind = "Main"

var MainRoutine = Actions{
	Info("Hello world"),
}

func init() {
	MainRoutine.Register(MainRoutineKind)
}
