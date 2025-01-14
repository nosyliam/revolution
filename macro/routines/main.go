package routines

import (
	"github.com/nosyliam/revolution/macro/routines/hive"
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const MainRoutineKind RoutineKind = "Main"

var MainRoutine = Actions{
	Condition(
		If(False(V[bool](RestartSleep))),
		Redirect(OpenRobloxRoutineKind),
	),
	Condition(
		If(Equal(MS[int]("counters.claimedHive"), -1)),
		Routine(hive.ClaimHiveRoutineKind),
	),
	Info("Idling")(Status),
	Sleep(1).Seconds(),
}

func init() {
	MainRoutine.Register(MainRoutineKind)
}
