package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/control"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	"github.com/pkg/errors"
)

const CheckRobloxRoutineKind RoutineKind = "CheckRoblox"

func openRoblox(macro *Macro) error {
	win, err := macro.WinManager.OpenRoblox()
	if err != nil {
		//_ = macro.Action(Sleep(1000))
		return RetrySignal
	}
	macro.Window = win
	return nil
}

func CheckRobloxRoutine(macro *Macro) []Action {
	return []Action{
		Condition(
			If(Nil(macro.Window)),
			openRoblox,
		),
		Condition(
			If(func() bool { return macro.Results.RetryCount > 10 }),
			func() error { return errors.New("failed to adjust and screenshot window") },
			If(NotNil(macro.Window.Screenshot())),
			func() error { macro.Results.RetryCount++; return RetrySignal },
		),
	}
}

func init() {
	control.Register(CheckRobloxRoutineKind, CheckRobloxRoutine)
}
