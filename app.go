package main

import (
	"context"
	"fmt"
	"github.com/nosyliam/revolution/macro"
	"github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/control"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
	"github.com/sqweek/dialog"
	"os"
)

type Macro struct {
	ctx context.Context

	config    *Object[Config]
	state     *Object[State]
	database  *Object[AccountDatabase]
	runtime   *Runtime
	eventBus  common.EventBus
	backend   common.Backend
	windowMgr *window.Manager

	interfaces map[string]*macro.Interface

	err chan string
}

func NewMacro(
	windowBackend window.Backend,
	controlBackend common.Backend,
) *Macro {
	return &Macro{
		windowMgr:  window.NewWindowManager(windowBackend),
		eventBus:   control.NewEventBus(controlBackend),
		backend:    controlBackend,
		interfaces: make(map[string]*macro.Interface),
	}
}

func (m *Macro) exitError(err error) {
	dialog.Message(err.Error()).Error()
	os.Exit(1)
}

func (m *Macro) startup(ctx context.Context) {
	AppContext = ctx
	m.runtime = NewRuntime(ctx)
	var err error
	m.config, err = NewConfig(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load configuration").Error()).Error()
		os.Exit(1)
	}
	m.state, err = NewState(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load state").Error()).Error()
		os.Exit(1)
	}
	m.database, err = NewDatabase(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load account database").Error()).Error()
		os.Exit(1)
	}

	presets := *Concrete[[]*Object[Settings]](m.config, "presets")
	var accounts = make(map[string]*Object[Settings])

	// Find the default preset for the default account (no join url; deep-link only)
	var defaultPresetName = *Concrete[string](m.state, "config.defaultPreset")
	if defaultPresetName == "" {
		defaultPresetName = "default"
	}
	if def := Concrete[*Object[Settings]](m.config, "presets[%s]", defaultPresetName); def != nil {
		accounts["default"] = *def
	} else {
		accounts["default"] = presets[0]
		if err := m.config.SetPath("defaultPreset", presets[0].Object().Name); err != nil {
			m.exitError(errors.Wrap(err, "failed to set default preset"))
		}
	}

	for _, account := range *Concrete[[]*Object[Account]](m.database, "accounts") {
		object := account.Object()
		if preset := Concrete[*Object[Settings]](m.config, "presets[%s]", object.Preset); preset != nil {
			accounts[object.Name] = *preset
		} else {
			accounts[object.Name] = presets[0]
			if err := account.SetPath("preset", presets[0].Object().Name); err != nil {
				m.exitError(errors.Wrap(err, fmt.Sprintf("failed to set preset for account %s", object.Name)))
			}
		}
	}

	for name, preset := range accounts {
		macroPath := fmt.Sprintf("macros[%s]", name)
		var macroState *Object[MacroState]
		if state := Concrete[*Object[MacroState]](m.state, macroPath); state == nil {
			if err = m.state.AppendPath(macroPath); err != nil {
				m.exitError(errors.Wrap(err, fmt.Sprintf("Failed to load state for macro %s", name)))
			}
			macroState = *Concrete[*Object[MacroState]](m.state, macroPath)
		} else {
			macroState = *state
		}
		m.interfaces[name] = macro.NewInterface(
			name,
			preset,
			macroState,
			m.database,
			m.windowMgr,
			m.eventBus,
			m.backend,
		)
	}
}

func (m *Macro) Start(instance string) {
	account := m.interfaces[instance]
	if *Concrete[bool](m.state, "macros[%s].paused", instance) {
		account.Unpause()
		return
	}
	account.Start()
}

func (m *Macro) Pause(instance string) {
	account := m.interfaces[instance]
	account.Pause()
}

func (m *Macro) Stop(instance string) {
	account := m.interfaces[instance]
	account.Stop()
}

func (m *Macro) StartAll() {
	for name, account := range m.interfaces {
		if !*Concrete[bool](m.state, "macros[%s].running", name) || *Concrete[bool](m.state, "macros[%s].paused", name) {
			account.Start()
		}
	}
}

func (m *Macro) PauseAll() {
	for name, account := range m.interfaces {
		if !*Concrete[bool](m.state, "macros[%s].paused", name) && *Concrete[bool](m.state, "macros[%s].running", name) {
			account.Pause()
		}
	}
}

func (m *Macro) StopAll() {
	for name, account := range m.interfaces {
		if *Concrete[bool](m.state, "macros[%s].running", name) {
			account.Stop()
		}
	}
}

func (m *Macro) SetAccountPreset(account, name string) string {
	ifc, ok := m.interfaces[account]
	if !ok {
		return fmt.Sprintf("Failed to find account \"%s\"", account)
	}
	if preset := Concrete[*Object[Settings]](m.config, "presets[%s]", name); preset != nil {
		if account != "default" {
			if err := m.database.SetPathf(name, "accounts[%s].preset", account); err != nil {
				return "Failed to set active preset"
			}
		} else {
			if err := m.state.SetPath("config.defaultPreset", name); err != nil {
				return "Failed to set active preset"
			}
		}
		ifc.Settings = *preset
	} else {
		return fmt.Sprintf("Failed to find preset \"%s\"", account)
	}
	return ""
}

func (m *Macro) DeleteAccount(name string) string {
	ifc, ok := m.interfaces[name]
	if !ok {
		return fmt.Sprintf("Failed to find account \"%s\"", name)
	}
	account := Concrete[Account](m.database, "accounts[%s]", name)
	if account == nil {
		return fmt.Sprintf("Failed to find account \"%s\"", name)
	}
	state := Concrete[MacroState](m.state, "macros[%s]", name)
	if state == nil {
		return fmt.Sprintf("Failed to find account \"%s\" in state", name)
	}
	if err := m.database.DeletePathf("accounts[%s]", name); err != nil {
		return fmt.Sprintf("Failed to delete account \"%s\" from database", name)
	}
	if err := m.database.DeletePathf("state[%s]", name); err != nil {
		return fmt.Sprintf("Failed to delete account \"%s\" from state", name)
	}
	account.Delete()
	if state.Paused {
		ifc.Unpause()
	}
	if state.Running {
		ifc.Stop()
	}
	if ifc.Macro.Window != nil {
		_ = ifc.Macro.Window.Close()
	}
	delete(m.interfaces, name)
	return ""
}

func (m *Macro) GetLoginCode() string {
	return ""
}

func (m *Macro) DeletePreset(name string) string {
	var preset *Object[Settings]
	presets := *Concrete[[]*Object[Settings]](m.config, "presets")
	for _, object := range presets {
		if object.Object().Name == name {
			preset = object
		}
	}
	if preset == nil {
		return fmt.Sprintf("Failed to find preset \"%s\"", preset)
	}

	for _, ifc := range m.interfaces {
		if ifc.Settings == preset {
			if ifc.Account != "default" {
				if err := m.database.SetPathf(presets[0].Object().Name, "accounts[%s].preset"); err != nil {
					return "Failed to set active preset"
				}
			} else {
				if err := m.state.SetPath("config.defaultPreset", presets[0].Object().Name); err != nil {
					return "Failed to set active preset"
				}
			}
			ifc.Settings = presets[0]
		}
	}
	return ""
}
