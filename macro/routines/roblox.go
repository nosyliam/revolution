package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const (
	OpenRobloxRoutineKind RoutineKind = "OpenRoblox"
	ResetRoutineKind      RoutineKind = "Reset"
)

func closeWindow(macro *Macro) error {
	if err := macro.Root.Window.Close(); err != nil {
		return err
	}

	macro.Root.Window = nil
	return nil
}

func openWindow(macro *Macro) error {
	if macro.Root.Window != nil {
		if err := closeWindow(macro); err != nil {
			return err
		}
	}
	ignoreLink := V[bool](UsePublicServer)(macro)
	if win, err := macro.WinManager.OpenWindow(macro.Account, macro.Database, macro.Settings, ignoreLink); err != nil {
		return err
	} else {
		macro.Root.Window = win
		return nil
	}
}

var LoadingImage = ImageSteps{
	SelectCoordinate(Change, 0, 30, 0, 150),
	Variance(4),
	Search("loading").Find(),
}

var ScienceImage = ImageSteps{
	SelectCoordinate(Change, 0, 30, 0, 150),
	Variance(4),
	Search("science").Find(),
}

var DisconnectImage = ImageSteps{
	SelectCoordinate(Change, 0, 30, 0, Sub(Height, 30)),
	Variance(2),
	Search("disconnected").Find(),
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
			Loop(
				For(180),
				Condition(
					If(Equal(Index(), 179)),
					Error("No BSS was found!")(Status, Discord),
					Sleep(5).Seconds(),
					Continue(1),
				),
				Condition(
					If(ExecError(TakeScreenshot)),
					Error("Failed to screenshot BSS!")(Status, Discord),
					Sleep(5).Seconds(),
					Continue(1),
					Else(),
					Condition(
						If(Image(LoadingImage...).Found()),
						Info("Game Open")(Status, Discord),
						Break(),
						If(Image(ScienceImage...).Found()),
						Info("Game Loaded")(Status, Discord),
						Terminate(),
						If(Image(DisconnectImage...).Found()),
						Info("Disconnected during reconnect")(Status, Discord),
						Sleep(5).Seconds(),
						Continue(1),
					),
				),
			),
			Loop(
				For(180),
				Condition(
					If(Equal(Index(), 179)),
					Error("BSS load timeout exceeded!")(Status, Discord),
					Sleep(5).Seconds(),
					Continue(1),
				),
				Condition(
					If(ExecError(TakeScreenshot)),
					Error("Failed to screenshot BSS!")(Status, Discord),
					Sleep(5).Seconds(),
					Continue(1),
					Else(),
					Condition(
						If(Or(
							Image(LoadingImage...).NotFound(),
							Image(ScienceImage...).Found(),
						)),
						Info("Game Loaded")(Status, Discord),
						Break(),
						If(Image(DisconnectImage...).Found()),
						Info("Disconnected during reconnect!")(Status, Discord),
						Sleep(5).Seconds(),
						Continue(1),
					),
				),
			),
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
			go macro.Scheduler.Start()
		}),
	),
	Set(RestartSleep, true),
}

func init() {
	OpenRobloxRoutine.Register(OpenRobloxRoutineKind)
}
