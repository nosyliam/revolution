package macro

import (
	"fmt"
	"github.com/nosyliam/revolution/macro/routines"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	"slices"
	"time"
)

const ClockTime = 50 * time.Millisecond

var intervals []*interval

type interval struct {
	priority int
	delayed  bool
	delay    int64
	kind     common.RoutineKind

	lastExec, enabled string
}

// Scheduler manages the execution of recurring tasks called interrupts. When an interrupt is activated, execution
// will unwind to the routine requesting the interrupt. An interrupt can exist in one of three forms:
//
// Immediate: Executed solely in the scheduler clock to immediately redirect
// execution in response to critical events (i.e. roblox closed, balloon blessing timer)
//
// Delayed: Executed at routine-specific points (i.e. gather interrupt).
// Interval interrupts may be registered as delayed interrupts.
//
// Interval: Executed at time-based intervals (i.e. every hour), based on the time since last execution.
// Examples of interval interrupts include planters, bug runs and wealth clock. Interval interrupts are
// only executed during the start of the main loop (after the hive backboard check) and are assigned a specific priority
type Scheduler struct {
	macro    *common.Macro
	close    chan struct{}
	redirect chan<- *common.RedirectExecution
	tick     int
}

func RegisterInterval(
	kind common.InterruptKind,
	routine common.RoutineKind,
	priority int,
	delay int64,
	lastExecutionPath, enabledPath string,
) {
	intervals = append(intervals, &interval{
		priority: priority,
		delay:    delay,
		delayed:  kind == common.DelayedInterrupt,
		kind:     routine,
		lastExec: lastExecutionPath,
		enabled:  enabledPath,
	})
}

func (s *Scheduler) Execute(interruptType common.InterruptKind) {
	var ivls []*interval

	for _, ivl := range intervals {
		if interruptType == common.DelayedInterrupt && ivl.delayed {
			ivls = append(ivls, ivl)
		} else if interruptType == common.IntervalInterrupt && !ivl.delayed {
			ivls = append(ivls, ivl)
		}
	}

	slices.SortFunc(ivls, func(a, b *interval) int {
		return a.priority - b.priority
	})

	for _, ivl := range ivls {
		if *config.Concrete[bool](s.macro.Settings, ivl.enabled) {
			var lastExec = *config.Concrete[int64](s.macro.State, ivl.lastExec)
			if time.Now().Unix() >= lastExec+ivl.delay*60 {
				s.macro.Routine(ivl.kind)
				_ = s.macro.State.SetPath(ivl.lastExec, time.Now().Unix())
			}
		}
	}
}

func (s *Scheduler) Close() {
	if s.close == nil {
		panic("scheduler already closed")
	}
	s.close <- struct{}{}
}

func (s *Scheduler) Tick() {
	defer func() {
		s.tick++
	}()
	// If we're not opening the window or unwinding a redirect, check the Roblox window
	if opening := s.macro.Scratch.Stack[0] == string(routines.OpenRobloxRoutineKind); !opening && !s.macro.Scratch.Redirect {
		if s.macro.Root.Window == nil {
			fmt.Println("failed 0")
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
			return
		}
		if err := s.macro.Root.Window.Fix(); err != nil {
			s.macro.Action(Error("Failed to adjust Roblox! Re-opening")(Status))
			s.macro.Action(Error("Failed to adjust Roblox: %s! Attempting to re-open", err)(Discord))
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
			return
		}
		if err := s.macro.Root.Window.TakeScreenshot(); err != nil {
			s.macro.Action(Error("Failed to screenshot Roblox! Re-opening")(Status))
			s.macro.Action(Error("Failed to screenshot Roblox: %s! Attempting to re-open", err)(Discord))
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
			return
		}
	} else if s.macro.Scratch.Redirect || opening {
		return
	}
}

func (s *Scheduler) Start() {
	s.close = make(chan struct{})
	for {
		select {
		case <-time.After(ClockTime):
			s.Tick()
		case <-s.close:
			s.close = nil
			return
		}
	}
}

func (s *Scheduler) Initialize(macro *common.Macro) {
	s.macro = macro
}

func NewScheduler(redirect chan<- *common.RedirectExecution) common.Scheduler {
	return &Scheduler{
		redirect: redirect,
	}
}

type InterruptAction struct {
	kind common.InterruptKind
}

func (a *InterruptAction) Execute(macro *common.Macro) error {
	macro.Scheduler.Execute(a.kind)
	return nil
}

func Interrupt() common.Action {
	return &InterruptAction{}
}
