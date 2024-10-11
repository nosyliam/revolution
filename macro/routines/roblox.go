package routines

import (
	"fmt"
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const (
	OpenRobloxRoutineKind RoutineKind = "OpenRoblox"
)

func closeWindow(macro *Macro) error {
	fmt.Println("window", macro.Window)
	if err := macro.Window.Close(); err != nil {
		return err
	}

	macro.Window = nil
	return nil
}

func openWindow(macro *Macro) error {
	ignoreLink := V[bool](UsePublicServer)(macro)
	if win, err := macro.WinManager.OpenWindow(macro.Account, macro.Database, macro.Settings, ignoreLink); err != nil {
		fmt.Println("open err", err)
		return err
	} else {
		macro.Window = win
		return nil
	}
}

func fallbackServerEnabled(macro *Macro) bool {
	return false
}

func privateServerAttempts(macro *Macro) int {
	return 0
}

var OpenRobloxRoutine = Actions{
	Set(RetryCount, 0),
	Set(UsePublicServer, false),
	Condition(
		If(NotNil(Window)),
		Info("Attempting to close Roblox")(Status, Discord),
		Loop(
			For(5),
			Condition(
				If(ExecError(closeWindow)),
				Error("Failed to close Roblox: %s! Attempt: %d", LastError, Index(0))(Status, Discord),
				Sleep(5).Seconds(),
				Else(),
				Break(),
			),
		),
	),
	Info("Opening Roblox")(Status, Discord),
	Loop(
		For(10),
		Condition(
			If(ExecError(openWindow)),
			Error("Failed to open Roblox! Attempt: %d", Index(0))(Status, Discord),
			Sleep(5).Seconds(),
			Condition(
				If(And(GreaterThanEq(V[int](RetryCount), 5), True(P[bool]("window.fallbackToPublicServer")))),
				Set(UsePublicServer, true),
				If(Equal(privateServerAttempts, 5)),
				Break(),
				Else(),
				Increment(RetryCount),
			),
			Else(),
			Break(),
		),
	),
	Condition(
		If(Nil(Window)),
		Error("Waiting 30 seconds before retrying")(Status, Discord),
		Sleep(30).Seconds(),
		Restart(),
	),
	Redirect(MainRoutineKind),
}

func init() {
	OpenRobloxRoutine.Register(OpenRobloxRoutineKind)
}
