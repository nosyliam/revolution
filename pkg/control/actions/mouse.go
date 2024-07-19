package actions

import (
	"github.com/nosyliam/revolution/pkg/control/common"
)

type mouseMoveAction struct {
	x, y int
}

func (a *mouseMoveAction) Execute(deps *common.Dependencies) error {
	<-deps.EventBus.MoveMouse(a.x, a.y)
	return nil
}

func MouseMoveAction(x, y int) common.Action {
	return &mouseMoveAction{x, y}
}
