package control

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type MouseMoveEvent struct {
	x, y int
}

func (e *MouseMoveEvent) Execute(backend common.Backend) {
	backend.MoveMouse(e.x, e.y)
}

type MouseScrollEvent struct {
	x, y int
}

func (e *MouseScrollEvent) Execute(backend common.Backend) {
	backend.ScrollMouse(e.x, e.y)
}

type KeyDownEvent struct {
	key common.Key
	id  int
}

func (e *KeyDownEvent) Execute(backend common.Backend) {
	backend.KeyDown(e.id, e.key)
}

type KeyUpEvent struct {
	key common.Key
	id  int
}

func (e *KeyUpEvent) Execute(backend common.Backend) {
	backend.KeyUp(e.id, e.key)
}
