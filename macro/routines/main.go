package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const MainRoutineKind RoutineKind = "Main"

var MainRoutine = Actions{
	Condition(
		If(False(V[bool](RestartSleep))),
		Redirect(OpenRobloxRoutineKind),
	),
	Set("test", 0),
	Loop(
		For(10),
		Condition(
			If(GreaterThan(V[int]("test"), 5)),
			Info("Second: %d", VI("test"))(Status),
			Else(),
			Info("First: %d", VI("test"))(Status),
		),
		Increment("test"),
		Sleep(1).Seconds(),
	),
	Sleep(1).Seconds(),
}

func init() {
	MainRoutine.Register(MainRoutineKind)
}
