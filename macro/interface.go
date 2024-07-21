package macro

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
)

type Interface struct {
	EventBus common.EventBus
	Backend  common.Backend
	Settings *config.Settings
	Logger   *logging.Logger

	Name string
}

func (m *Interface) Start() {
}

func (m *Interface) SetStatus() {

}

func NewMacroInterface(name string, settings *config.Settings, eventBus common.EventBus, backend common.Backend) *Interface {
	return &Interface{
		EventBus: eventBus,
		Backend:  backend,
		Settings: settings,
		Logger:   logging.NewLogger(name, settings),
		Name:     name,
	}
}
