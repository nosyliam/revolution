package control

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type event struct {
	event common.Event
	wait  chan<- struct{}
}

type eventBusImpl struct {
	queue       chan event
	attachedPid map[int]bool
	backend     common.Backend
}

func (e *eventBusImpl) KeyDown(macro *common.Macro, key common.Key) common.Receiver {
	if _, ok := e.attachedPid[macro.Root.Window.PID()]; ok {
		e.backend.AttachInput(macro.Root.Window.PID())
		e.attachedPid[macro.Root.Window.PID()] = true
	}
	ch := make(chan struct{})
	e.queue <- event{&KeyDownEvent{Event{macro.Root.Window}, key}, ch}
	return ch

}

func (e *eventBusImpl) KeyUp(macro *common.Macro, key common.Key) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&KeyUpEvent{Event{macro.Root.Window}, key}, ch}
	return ch
}

func (e *eventBusImpl) MoveMouse(macro *common.Macro, x, y int) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&MouseMoveEvent{Event{macro.Root.Window}, x, y}, ch}
	return ch
}

func (e *eventBusImpl) ScrollMouse(macro *common.Macro, x, y int) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&MouseScrollEvent{Event{macro.Root.Window}, x, y}, ch}
	return ch
}

func (e *eventBusImpl) Start() {
	for {
		evt, ok := <-e.queue
		if !ok {
			break
		}
		evt.event.Execute(e.backend)
		evt.wait <- struct{}{}
	}
}

func NewEventBus(backend common.Backend) common.EventBus {
	return &eventBusImpl{backend: backend, queue: make(chan event, 1), attachedPid: make(map[int]bool)}
}
