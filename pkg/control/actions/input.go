package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type mouseMoveAction struct {
	x, y int
}

func (a *mouseMoveAction) Execute(macro *common.Macro) error {
	<-macro.EventBus.MoveMouse(macro, a.x, a.y)
	return nil
}

func MouseMoveAction(x, y int) common.Action {
	return &mouseMoveAction{x, y}
}
