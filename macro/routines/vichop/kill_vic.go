package vichop

import (
	"github.com/nosyliam/revolution/macro/routines/hive"
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const KillVicRoutineKind RoutineKind = "kill-vic"

var KillVic = Actions{
	Info("Killing Vicious Bee: %s", V[string](VicField))(Status, Discord),
	Routine(hive.ClaimHiveRoutineKind),
	Sleep(10).Seconds(),
	Set(VicField, ""),
	Condition(
		If(True(P[bool]("vicHop.serverHop"))),
		Redirect("OpenRoblox"),
	),
}

func init() {
	KillVic.Register(KillVicRoutineKind)
}
