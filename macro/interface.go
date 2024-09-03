package macro

import (
	"github.com/nosyliam/revolution/macro/routines"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/control"
	"github.com/nosyliam/revolution/pkg/logging"
)

type Interface struct {
	EventBus common.EventBus
	Backend  common.Backend
	Settings config.Reactive
	Database config.Reactive
	State    config.Reactive
	Logger   *logging.Logger
	Name     string

	pause    chan struct{}
	stop     chan struct{}
	redirect chan *common.RedirectExecution
}

func (i *Interface) Start() {
	macro := &common.Macro{
		EventBus: i.EventBus,
		Backend:  i.Backend,
		Settings: i.Settings,
		Database: i.Database,
		Logger:   i.Logger,
		State:    i.State,
		Results:  &common.ActionResults{},
	}

	var unpause chan struct{}
	stop := make(chan struct{}, 1)
	pause := make(chan (<-chan struct{}), 1)
	err := make(chan string, 1)

	go func() {
		for {
			select {
			case errStr := <-err:
				i.SendError(errStr)
				i.Pause()
			case <-i.pause:
				if unpause != nil {
					unpause <- struct{}{}
					unpause = nil
				}
				unpause = make(chan struct{})
				pause <- unpause
			case <-i.stop:
				stop <- struct{}{}
				break
			}
		}
	}()

	main := common.Routines[routines.MainRoutineKind]
	control.ExecuteRoutine(macro, main, stop, pause, make(chan<- string), err, i.redirect)
}

func (i *Interface) SendError(err string) {
	// TODO: send error to UI
}

func (i *Interface) Stop() {
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

func (i *Interface) RegisterEventListeners() {

}

func NewInterface(
	name string,
	state config.Reactive,
	settings config.Reactive,
	database config.Reactive,
	eventBus common.EventBus,
	backend common.Backend,
) *Interface {
	return &Interface{
		EventBus: eventBus,
		Backend:  backend,
		Settings: settings,
		Database: database,
		Logger:   logging.NewLogger(name, settings),
		State:    state,
		Name:     name,

		pause:    make(chan struct{}),
		stop:     make(chan struct{}),
		redirect: make(chan *common.RedirectExecution, 1),
	}
}
