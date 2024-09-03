package main

import (
	"context"
	"github.com/nosyliam/revolution/macro"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
	"github.com/sqweek/dialog"
	"os"
)

type Macro struct {
	ctx context.Context

	config     config.Reactive
	state      config.Reactive
	database   config.Reactive
	runtime    *config.Runtime
	eventBus   common.EventBus
	windowMgr  *window.Manager
	interfaces []*macro.Interface

	activeAccount int
	err           chan string
}

func NewMacro(
	windowMgr *window.Manager,
	eventBus common.EventBus,
) *Macro {
	return &Macro{windowMgr: windowMgr, eventBus: eventBus}
}

func (m *Macro) loadInstance() {

}

func (m *Macro) startup(ctx context.Context) {
	config.AppContext = ctx
	m.runtime = config.NewRuntime(ctx)
	var err error
	m.config, err = config.NewConfig(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load configuration").Error()).Error()
		os.Exit(1)
	}
	m.state, err = config.NewState(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load state").Error()).Error()
		os.Exit(1)
	}
	m.database, err = config.NewDatabase(m.runtime)
	if err != nil {
		dialog.Message(errors.Wrap(err, "Failed to load account database").Error()).Error()
		os.Exit(1)
	}
}

func (m *Macro) LoadSettings(preset string) *config.Settings {
	return nil
}

func (m *Macro) LoadState(account string) *config.MacroState {
	return nil
}

func (m *Macro) SetActiveInstance() {

}

func (m *Macro) Pause() {

}

func (m *Macro) Resume() {

}

func (m *Macro) PauseAll() {

}

func (m *Macro) ResumeAll() {

}
