package vichop

import (
	"github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/networking"
	"github.com/pkg/errors"
)

type Manager struct {
	Dataset *Dataset
	presets map[string]string
	states  map[string]*Object[MacroState]
	macros  map[string]*common.Macro

	settings *Object[Config]
	state    *Object[State]
}

func NewManager(settings *Object[Config], state *Object[State]) *Manager {
	return &Manager{
		Dataset: NewDataset(state),

		state:    state,
		settings: settings,

		presets: make(map[string]string),
		states:  make(map[string]*Object[MacroState]),
		macros:  make(map[string]*common.Macro),
	}
}

func (m *Manager) HandleRoleChange(account string, role string) {
	if macro, ok := m.macros[account]; ok {
		macro.Network.Client.SetRole(common.ClientRole(role))
	}
}

func (m *Manager) RegisterPreset(preset *Object[Settings]) {
	if err := preset.ListenPath("vicHop.role", func(_ ListenOp, value interface{}) {
		if !preset.Object().VicHop.Object().Enabled {
			return
		}
		for account, activePreset := range m.presets {
			if preset.Object().Name == activePreset {
				m.HandleRoleChange(account, value.(string))
			}
		}
	}); err != nil {
		panic(err)
	}
	if err := preset.ListenPath("vicHop.enabled", func(_ ListenOp, value interface{}) {
		if value.(bool) {
			return
		}
		for account, activePreset := range m.presets {
			if preset.Object().Name == activePreset {
				m.HandleRoleChange(account, common.InactiveClientRole)
			}
		}
	}); err != nil {
		panic(err)
	}
}

func (m *Manager) RegisterState(state *Object[MacroState]) {
	m.states[state.Object().AccountName] = state
}

func (m *Manager) UnregisterState(state *Object[MacroState]) {
	delete(m.states, state.Object().AccountName)
}

func (m *Manager) Start() {
	m.settings.Object().Presets.ForEachObject(func(value *Object[Settings]) {
		m.RegisterPreset(value)
	})
	m.state.Object().Macros.ForEachObject(func(value *Object[MacroState]) {
		m.RegisterState(value)
	})
	if err := m.settings.ListenPath("presets", func(op ListenOp, value interface{}) {
		if op == Append {
			m.RegisterPreset(value.(*Object[Settings]))
		}
	}); err != nil {
		panic(err)
	}
	if err := m.state.ListenPath("macros", func(op ListenOp, value interface{}) {
		if op == Append {
			m.RegisterState(value.(*Object[MacroState]))
		} else {
			m.UnregisterState(value.(*Object[MacroState]))
		}
	}); err != nil {
		panic(err)
	}
	if err := m.state.ListenPath("config.defaultPreset", func(_ ListenOp, value interface{}) {
		m.presets["Default"] = value.(string)
		role, _ := m.settings.GetPathf("presets[%s].vicHop.role", value.(string))
		m.HandleRoleChange("Default", role.(string))
	}); err != nil {
		panic(err)
	}
}

func (m *Manager) RegisterMacro(macro *common.Macro) error {
	m.macros[macro.MacroState.Object().AccountName] = macro
	if settings := macro.Settings.Object().VicHop.Object(); settings.Enabled {
		if err := macro.Network.Client.SetRole(common.ClientRole(settings.Role)); errors.Is(err, networking.InactiveNetworkError) {
			return errors.New("please connect to a network relay to use Vic Hop")
		} else if err != nil {
			return err
		}
	}
	return nil
}
