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
	if err := macro.Root.Window.Close(); err != nil {
		return err
	}

	macro.Root.Window = nil
	return nil
}

func openWindow(macro *Macro) error {
	ignoreLink := V[bool](UsePublicServer)(macro)
	if win, err := macro.WinManager.OpenWindow(macro.Account, macro.Database, macro.Settings, ignoreLink); err != nil {
		fmt.Println("open err", err)
		return err
	} else {
		macro.Root.Window = win
		return nil
	}
}

var OpenRobloxRoutine = Actions{
	Condition(If(True(V[bool](RestartSleep))), Sleep(5).Seconds()),
	Set(RetryCount, 0),
	Set(UsePublicServer, false),
	Condition(
		If(NotNil(Window)),
		Info("Attempting to close Roblox")(Status, Discord),
		Logic(closeWindow),
		Sleep(3).Seconds(),
	),
	Info("Opening Roblox")(Status, Discord),
	Loop(
		For(1, 11),
		Logic(func() {
			fmt.Println("executing index")
		}),
		Condition(
			If(ExecError(openWindow)),
			Error("Failed to open Roblox! Attempt: %d", Index(0))(Status, Discord),
			Sleep(5).Seconds(),
			Condition(
				If(And(GreaterThanEq(Index(), 5), True(P[bool]("window.fallbackToPublicServer")))),
				Set(UsePublicServer, true),
				If(Equal(Index(), 5)),
				Break(),
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
	Condition(
		If(False(V[bool](RestartSleep))),
		Logic(func(macro *Macro) {
			fmt.Println("start scheduler")
			go macro.Scheduler.Start()
		}),
	),
	Set(RestartSleep, true),
}

func init() {
	OpenRobloxRoutine.Register(OpenRobloxRoutineKind)
}
