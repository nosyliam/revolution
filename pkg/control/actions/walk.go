package actions

import "github.com/nosyliam/revolution/pkg/common"

type SleepAction struct {
	ms      int
	seconds bool
}

func (a *SleepAction) Seconds() *SleepAction {
	a.seconds = true
	return a
}

func (a *SleepAction) Execute(macro *common.Macro) error {
	var time = a.ms
	if a.seconds {
		time *= 1000
	}
	var interrupt = make(chan struct{})
	macro.Backend.Sleep(time, interrupt)
	return nil
}

func Sleep(ms int) *SleepAction {
	return &SleepAction{ms, false}
}
