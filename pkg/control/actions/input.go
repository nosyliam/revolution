package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/movement"
)

type mouseMoveAction struct {
	x, y int
}

func (a *mouseMoveAction) Execute(macro *common.Macro) error {
	<-macro.EventBus.MoveMouse(macro, a.x, a.y)
	return nil
}

func MoveMouse(x, y int) common.Action {
	return &mouseMoveAction{x, y}
}

type keyPressAction struct {
	key common.Key
}

func (a *keyPressAction) Execute(macro *common.Macro) error {
	<-macro.EventBus.KeyDown(macro, a.key)
	movement.Sleep(50, macro) // TODO: Key delay
	<-macro.EventBus.KeyUp(macro, a.key)
	return nil
}

func KeyPress(key common.Key) common.Action {
	return &keyPressAction{key}
}

type keyDownAction struct {
	key common.Key
}

func (a *keyDownAction) Execute(macro *common.Macro) error {
	<-macro.EventBus.KeyDown(macro, a.key)
	return nil
}

func KeyDown(key common.Key) common.Action {
	return &keyDownAction{key}
}

type keyUpAction struct {
	key common.Key
}

func (a *keyUpAction) Execute(macro *common.Macro) error {
	<-macro.EventBus.KeyUp(macro, a.key)
	return nil
}

func KeyUp(key common.Key) common.Action {
	return &keyUpAction{key}
}
