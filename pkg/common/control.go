package common

import (
	revimg "github.com/nosyliam/revolution/pkg/image"
	"image"
)

type Key int
type Receiver <-chan struct{}

var Send = struct{}{}

const (
	Forward Key = iota
	Backward
	Left
	Right
	RotLeft
	RotRight
	RotUp
	RotDown
	ZoomIn
	ZoomOut
	E
	R
	L
	Esc
	Enter
	LShift
	Space
	One
	Two
	Three
	Four
	Five
	Six
	Seven
)

type InterruptKind int

const (
	DelayedInterrupt InterruptKind = iota
	IntervalInterrupt
)

type Backend interface {
	KeyDown(pid int, key Key)
	KeyUp(pid int, key Key)
	MoveMouse(x, y int)
	ScrollMouse(x, y int)
	Sleep(ms int, interrupt Receiver)
	SleepAsync(ms int, interrupt Receiver) Receiver
}

type EventBus interface {
	Start()
	KeyDown(pid int, key Key) Receiver
	KeyUp(pid int, key Key) Receiver
	MoveMouse(x, y int) Receiver
	ScrollMouse(x, y int) Receiver
}

type Scheduler interface {
	Execute(interruptType InterruptKind)
	Start()
	Close()
	Initialize(macro *Macro)
}

type Event interface {
	Execute(backend Backend)
}

type BuffDetector interface {
	Tick(origin *revimg.Point, screenshot *image.RGBA)
	MoveSpeed() float64
	Watch() chan struct{}
}

type PatternLoader interface {
	Start() error
	Patterns() []string
	Exists(pattern string) bool
	Execute(macro *Macro, pattern string) error
}
