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
	Settings *config.Settings
	State    *config.MacroState
	Logger   *logging.Logger
	Name     string

	pause  chan struct{}
	stop   chan struct{}
	status chan string
}

func (i *Interface) Start() {
	macro := &common.Macro{
		EventBus: i.EventBus,
		Backend:  i.Backend,
		Settings: i.Settings,
		Logger:   i.Logger,
		State:    i.State,
		Results:  &common.ActionResults{},
	}

	var unpause chan struct{}
	stop := make(chan struct{})
	pause := make(chan (<-chan struct{}))
	err := make(chan string)
	status := make(chan string)

	go func() {
		for {
			select {
			case statStr := <-status:
				i.SendStatus(statStr)
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

	main := control.Routines[routines.MainRoutineKind]
	control.ExecuteRoutine(macro, main(macro), stop, pause, status, err)
}

func (i *Interface) SendError(err string) {
	// TODO: send error to UI
}

func (i *Interface) SendStatus(status string) {
	// TODO: send status to UI
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
	state *config.State,
	settings *config.Settings,
	eventBus common.EventBus,
	backend common.Backend,
) *Interface {
	return &Interface{
		EventBus: eventBus,
		Backend:  backend,
		Settings: settings,
		Logger:   logging.NewLogger(name, settings),
		State:    state.State(name),
		Name:     name,

		pause:  make(chan struct{}),
		stop:   make(chan struct{}),
		status: make(chan string),
	}
}
