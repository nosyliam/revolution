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
	State    *config.Object[config.State]
	Settings *config.Object[config.Settings]
	Database *config.Object[config.AccountDatabase]
	WinMgr   *window.Manager
	Logger   *logging.Logger
	Account  string

	pause    chan struct{}
	unpause  chan struct{}
	stop     chan struct{}
	quit     chan struct{}
	redirect chan *common.RedirectExecution
}

func (i *Interface) Start() {
	macro := &common.Macro{
		Account:    i.Account,
		EventBus:   i.EventBus,
		Backend:    i.Backend,
		Settings:   i.Settings,
		Database:   i.Database,
		Logger:     i.Logger,
		State:      i.State,
		WinManager: i.WinMgr,
		Results:    &common.ActionResults{},
	}

	var unpause chan struct{}
	stop := make(chan struct{}, 1)
	pause := make(chan (<-chan struct{}), 1)
	err := make(chan string, 1)

	_ = i.State.SetPathf(true, "macros[%s].running", i.Account)

	go func() {
		for {
			select {
			case errStr := <-err:
				i.SendError(errStr)
				i.Pause()
			case <-i.pause:
				if i.unpause != nil {
					_ = i.State.SetPathf(false, "macros[%s].paused", i.Account)
					i.unpause <- struct{}{}
					i.unpause = nil
				}
				unpause = make(chan struct{})
				pause <- unpause
			case <-i.stop:
				_ = i.State.SetPathf(false, "macros[%s].running", i.Account)
				stop <- struct{}{}
				break
			case <-i.quit:
				_ = i.State.SetPathf(false, "macros[%s].running", i.Account)
				stop <- struct{}{}
				_ = macro.Window.Close()
				break
			}
		}
	}()

	main := common.Routines[routines.MainRoutineKind]
	control.ExecuteRoutine(macro, main, stop, pause, make(chan<- string), err, i.redirect)
}

func (i *Interface) SendError(err string) {
	dialog.Message(err).Error()
}

func (i *Interface) Stop() {
	if i.unpause != nil {
		i.unpause <- struct{}{}
	}
	if len(i.stop) != 0 {
		return
	}
	i.stop <- struct{}{}
}

func (i *Interface) Pause() {
	_ = i.State.SetPathf(true, "macros[%s].paused", i.Account)
	if len(i.pause) != 0 {
		return
	}
	i.pause <- struct{}{}
}

func (i *Interface) Unpause() {
	if i.unpause == nil || len(i.unpause) != 0 {
		return
	}
	i.unpause <- struct{}{}
}

func NewInterface(
	account string,
	settings *config.Object[config.Settings],
	state *config.Object[config.State],
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

		pause:    make(chan struct{}),
		stop:     make(chan struct{}),
		redirect: make(chan *common.RedirectExecution, 1),
	}
}
