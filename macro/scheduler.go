package macro

import (
	"github.com/nosyliam/revolution/pkg/common"
	"time"
)

const ClockTime = 50 * time.Millisecond

type Interrupt struct {
	routine  common.RoutineKind
	priority int
}

type Scheduler struct {
	macro      *common.Macro
	close      chan struct{}
	err        chan<- string
	status     chan<- string
	immediate  []Interrupt
	delayed    []Interrupt
	interval   []Interrupt
	retryCount int
}

func (s *Scheduler) Execute(macro *common.Macro, interruptType common.InterruptType) error {
	var interrupts []Interrupt
	switch interruptType {
	case common.ImmediateInterrupt:
		interrupts = s.immediate
	case common.DelayedInterrupt:
		interrupts = s.delayed
	case common.IntervalInterrupt:
		interrupts = s.interval
	}

	if len(interrupts) == 0 {
		return nil
	}

}

func (s *Scheduler) Close() {
	if s.close == nil {
		panic("scheduler already closed")
	}
	s.close <- struct{}{}
}

func (s *Scheduler) Tick() error {
	if s.macro.Window == nil {

	}
}

func (s *Scheduler) Start() {
	s.close = make(chan struct{})
	for {
		select {
		case <-time.After(ClockTime):
			if err := s.Tick(); err != nil {
				s.err <- err.Error()
			}
		case <-s.close:
			s.close = nil
		}
	}
}

func NewScheduler(macro *common.Macro, err chan<- string, status chan<- string) common.Scheduler {
	return &Scheduler{
		macro:  macro,
		err:    err,
		status: status,
	}
}

type immediateInterruptAction struct{}

func (a *immediateInterruptAction) Execute(macro *common.Macro) error {

}

func Immediate() common.Action {
	return &immediateInterruptAction{}
}
