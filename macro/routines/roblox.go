package routines

import "github.com/nosyliam/revolution/pkg/control"

const CheckRobloxRoutine control.RoutineKind = "CheckRoblox"

func init() {
	control.NewRoutine(CheckRobloxRoutine)
}
