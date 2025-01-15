package vichop

import (
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
)

const VicSearchRoutineKind RoutineKind = "VicSearch"

var VicSearch = Actions{
	Condition(
		If(False(P[bool]("vicHop.enabled"))),
		Terminate(),
	),
	// Check if we're being called from OpenRoblox so that we can find a new server
	Condition(
		If(And(Equal(V[string](GameInstance), ""), True(P[bool]("vicHop.serverHop")))),
		Info("Finding new server")(Status, Discord),
		Set(GameInstance, func(macro *Macro) string {
			server, err := macro.VicHop.FindServer(macro)
			if server == "" {
				panic("no server")
			}
			if err != nil {
				macro.SetError(err, "Failed to find server")
			}
			return server
		}),
		Set(HopServer, true),
		Terminate(),
		If(Equal(V[string](GameInstance), "")),
		Set(HopServer, true),
		Info("Waiting for server")(Status, Discord),
		Loop(Forever(), Sleep(100)),
	),
	Set(VicField, ""),
	// Check for night
	Routine(DetectNightRoutineKind),
	Condition(
		If(False(V[bool](NightDetected))),
		Redirect("OpenRoblox"),
	),
	// Perform the vic search. Vic manager will handle redirecting on detection
	// We'll continue to scan for vicious bee attacking/defeated GUIs in case another player finds it
	ExecutePattern("vic_path").Async(),
	Loop(Until(False(PatternExecuting)), Sleep(100)),
	Redirect("OpenRoblox"),
}

func init() {
	VicSearch.Register(VicSearchRoutineKind)
}
