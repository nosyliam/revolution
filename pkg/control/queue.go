package control

import "github.com/nosyliam/revolution/pkg/control/common"

type event struct {
	event common.Event
	wait  chan<- struct{}
}

type eventBusImpl struct {
	queue   chan event
	backend common.Backend
}

func (e *eventBusImpl) KeyDown(pid int, key common.Key) <-chan struct{} {
	ch := make(chan struct{})

}

func (e *eventBusImpl) KeyUp(pid int, key common.Key) <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (e *eventBusImpl) MoveMouse(x, y int) <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (e *eventBusImpl) ScrollMouse(x, y int) <-chan struct{} {
	//TODO implement me
	panic("implement me")
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
