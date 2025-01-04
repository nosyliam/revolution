package control

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type event struct {
	event common.Event
	wait  chan<- struct{}
}

type eventBusImpl struct {
	queue   chan event
	backend common.Backend
}

func (e *eventBusImpl) KeyDown(pid int, key common.Key) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&KeyDownEvent{Event{pid}, key}, ch}
	return ch

}

func (e *eventBusImpl) KeyUp(pid int, key common.Key) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&KeyUpEvent{Event{pid}, key}, ch}
	return ch
}

func (e *eventBusImpl) MoveMouse(x, y int) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&MouseMoveEvent{x, y}, ch}
	return ch
}

func (e *eventBusImpl) ScrollMouse(x, y int) common.Receiver {
	ch := make(chan struct{})
	e.queue <- event{&MouseScrollEvent{x, y}, ch}
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

func (e *eventBusImpl) Close() {
	close(e.queue)
}

func NewEventBus(backend common.Backend) common.EventBus {
	return &eventBusImpl{backend: backend}
}
