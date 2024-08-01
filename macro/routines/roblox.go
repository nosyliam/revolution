package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/control"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const CheckRobloxRoutineKind RoutineKind = "CheckRoblox"

func closeWindow(macro *Macro) error {
	return macro.Window.Close()
}

var CheckRobloxRoutine = Actions{
	Condition(
		If(NotNil(Window)),
		Status("Attempting to close Roblox"),
		Loop(
			For(5),
			Condition(
				If(ExecError(closeWindow)),
				Error("Failed to close Roblox! Attempt %d", Index(0)),
				Else(),
				Break(),
			),
		),
	),
	Status("Opening Roblox"),
}

func init() {
	control.Register(CheckRobloxRoutineKind, CheckRobloxRoutine)
}
