package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type mouseMoveAction struct {
	x, y int
}

func (a *mouseMoveAction) Execute(deps *common.Macro) error {
	<-deps.EventBus.MoveMouse(a.x, a.y)
	return nil
}

func MouseMoveAction(x, y int) common.Action {
	return &mouseMoveAction{x, y}
}
