package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const (
	OpenRobloxRoutineKind RoutineKind = "OpenRoblox"
)

func closeWindow(macro *Macro) error {
	if err := macro.Window.Close(); err != nil {
		return err
	}

	macro.Window = nil
	return nil
}

func openWindow(macro *Macro) error {
	return nil
}

func fallbackServerEnabled(macro *Macro) bool {
	return macro.Settings.FallbackToPublicServer
}

func privateServerAttempts(macro *Macro) int {
	return macro.State.PrivateServerAttempts
}

var OpenRobloxRoutine = Actions{
	Condition(
		If(NotNil(Window)),
		Info("Attempting to close Roblox").Status().Discord(),
		Loop(
			For(5),
			Condition(
				If(ExecError(closeWindow)),
				Error("Failed to close Roblox: %s! Attempt %d", LastError, Index(0)).Status().Discord(),
				Sleep(5).Seconds(),
				Else(),
				Break(),
			),
		),
	),
	Info("Opening Roblox").Status().Discord(),
	Loop(
		For(10),
		Condition(
			If(ExecError(openWindow)),
			Error("Failed to open Roblox: %s! Attempt %d", LastError, Index(0)).Status().Discord(),
			Sleep(5).Seconds(),
			Condition(
				If(And(Equal(privateServerAttempts, 5), True(fallbackServerEnabled))),
				Logic(func(macro *Macro) {
					macro.State.UsePublicServer = true
				}),
				If(Equal(privateServerAttempts, 5)),
				Break(),
				Else(),
				Logic(func(macro *Macro) {
					macro.State.PrivateServerAttempts++
				}),
			),
			Else(),
			Break(),
		),
	),
	Condition(
		If(Nil(Window)),
		Error("Waiting 30 seconds before retrying").Status().Discord(),
		Sleep(30).Seconds(),
		Restart(),
	),
	Logic(func(macro *Macro) {
		macro.State.PrivateServerAttempts = 0
		macro.State.UsePublicServer = false
	}),
	Redirect(MainRoutineKind),
}

func init() {
	OpenRobloxRoutine.Register(OpenRobloxRoutineKind)
}
