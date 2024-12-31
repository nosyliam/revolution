package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	time "time"
)

type SleepAction struct {
	ms      int
	seconds bool
}

func (a *SleepAction) Seconds() *SleepAction {
	a.seconds = true
	return a
}

func (a *SleepAction) Execute(macro *common.Macro) error {
	var ms = a.ms
	if a.seconds {
		ms *= 1000
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}

func Sleep(ms int) *SleepAction {
	return &SleepAction{ms, false}
}
