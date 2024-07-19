package control

import "github.com/nosyliam/revolution/pkg/control/common"

type MouseMoveEvent struct {
	x, y int
}

func (e *MouseMoveEvent) Execute(_ int, backend common.Backend) {
	backend.MoveMouse(e.x, e.y)
}

type KeyDownEvent struct {
	key common.Key
}

func (e *KeyDownEvent) Execute(id int, backend common.Backend) {

	backend.KeyDown(id, e.key)
}
