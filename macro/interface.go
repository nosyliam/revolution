package macro

import (
	"github.com/nosyliam/revolution/macro/routines"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/control"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/sqweek/dialog"
)

type Interface struct {
	EventBus common.EventBus
	Backend  common.Backend
	State    *config.Object[config.MacroState]
	Settings *config.Object[config.Settings]
	Database *config.Object[config.AccountDatabase]
	WinMgr   *window.Manager
	Logger   *logging.Logger
	Macro    *common.Macro
	Account  string

	pause    chan struct{}
	unpause  chan struct{}
	stop     chan struct{}
	quit     chan struct{}
	redirect chan *common.RedirectExecution
}

func (i *Interface) Start() {
	pause := make(chan (<-chan struct{}), 1)
	i.Macro = &common.Macro{
		Account:    i.Account,
		EventBus:   i.EventBus,
		Backend:    i.Backend,
		Settings:   i.Settings,
		Database:   i.Database,
		Logger:     i.Logger,
		WinManager: i.WinMgr,
		Scratch:    config.NewScratch(),
		Results:    &common.ActionResults{},
		Pause:      pause,
	}

	stop := make(chan struct{}, 1)
	err := make(chan string, 1)
	status := make(chan string)

	_ = i.State.SetPath("running", true)

	go func() {
		for {
			select {
			case stat := <-status:
				_ = i.State.SetPath("status", stat)
			case errStr := <-err:
				i.SendError(errStr)
				i.Pause()
			case <-i.pause:
				if i.unpause != nil {
					_ = i.State.SetPath("paused", false)
					i.unpause <- struct{}{}
					i.unpause = nil
					continue
				}
				_ = i.State.SetPath("paused", true)
				i.unpause = make(chan struct{}, 1)
				pause <- i.unpause
			case <-i.stop:
				_ = i.State.SetPath("running", false)
				stop <- struct{}{}
				return
			}
		}
	}()

	main := common.Routines[routines.MainRoutineKind]
	control.ExecuteRoutine(i.Macro, main, stop, status, err, i.redirect)
}

func (i *Interface) SendError(err string) {
	dialog.Message(err).Error()
}

func (i *Interface) Stop() {
	if i.unpause != nil {
		i.unpause <- struct{}{}
		i.unpause = nil
		_ = i.State.SetPath("paused", false)
	}
	if len(i.stop) != 0 {
		return
	}
	i.stop <- struct{}{}
}

func (i *Interface) Pause() {
	if len(i.pause) != 0 {
		return
	}
	i.pause <- struct{}{}
}

func (i *Interface) Unpause() {
	if i.unpause == nil || len(i.unpause) != 0 {
		return
	}
	i.pause <- struct{}{}
}

func NewInterface(
	account string,
	settings *config.Object[config.Settings],
	state *config.Object[config.MacroState],
	database *config.Object[config.AccountDatabase],
	winMgr *window.Manager,
	eventBus common.EventBus,
	backend common.Backend,
) *Interface {
	return &Interface{
		EventBus: eventBus,
		Backend:  backend,
		Settings: settings,
		Database: database,
		Logger:   logging.NewLogger(account, settings),
		State:    state,
		WinMgr:   winMgr,
		Account:  account,

		pause:    make(chan struct{}, 1),
		stop:     make(chan struct{}, 1),
		redirect: make(chan *common.RedirectExecution, 1),
	}
}
