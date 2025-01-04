package control

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type Event struct {
	id int
}

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
	Event
	key common.Key
}

func (e *KeyDownEvent) Execute(backend common.Backend) {
	backend.KeyDown(e.id, e.key)
}

type KeyUpEvent struct {
	Event
	key common.Key
}

func (e *KeyUpEvent) Execute(backend common.Backend) {
	backend.KeyUp(e.id, e.key)
}
