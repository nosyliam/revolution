package macro

import (
	"github.com/nosyliam/revolution/macro/routines"
	"github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	"time"
)

const ClockTime = 50 * time.Millisecond

type InterruptMap map[common.RoutineKind]int

type interval struct {
	priority    int
	delayed     bool
	delay       time.Duration
	kind        common.RoutineKind
	getLastExec func() time.Time
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
// Interval: Executed at time-based intervals (i.e. every hour), based on a cron or time since last execution.
// Examples of interval interrupts include planters, bug runs and wealth clock.
type Scheduler struct {
	macro     *common.Macro
	close     chan struct{}
	redirect  chan<- *common.RedirectExecution
	delayed   InterruptMap
	interval  InterruptMap
	intervals []*interval
}

func (s *Scheduler) Execute(interruptType common.InterruptType) error {
	var interrupts InterruptMap
	switch interruptType {
	case common.DelayedInterrupt:
		interrupts = s.delayed
	case common.IntervalInterrupt:
		interrupts = s.interval
	}

	if len(interrupts) == 0 {
		return nil
	}
	return nil
}

func (s *Scheduler) Close() {
	if s.close == nil {
		panic("scheduler already closed")
	}
	s.close <- struct{}{}
}

func (s *Scheduler) Tick() {
	// If we're not opening the window or unwinding a redirect, check the Roblox window
	if opening := s.macro.State.Stack[0] == routines.OpenRobloxRoutineKind; !opening && len(s.redirect) == 0 {
		if s.macro.Window == nil {
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
		}
		if err := s.macro.Window.Fix(); err != nil {
			s.macro.Action(Error("Failed to adjust Roblox: %s! Attempting to close and re-open", LastError).
				Status().Discord())
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
			return
		}
		if err := s.macro.Window.Screenshot(); err != nil {
			s.macro.Action(Error("Failed to screenshot Roblox: %s! Attempting to close and re-open", LastError).
				Status().Discord())
			s.redirect <- &common.RedirectExecution{Routine: routines.OpenRobloxRoutineKind}
			return
		}
	} else if len(s.redirect) > 0 || opening {
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

func NewScheduler(macro *common.Macro, redirect chan<- *common.RedirectExecution) common.Scheduler {
	return &Scheduler{
		macro:    macro,
		redirect: redirect,
	}
}

type immediateInterruptAction struct{}

func (a *immediateInterruptAction) Execute(macro *common.Macro) error {

}

func IntervalInterrupt() common.Action {
	return &immediateInterruptAction{}
}
