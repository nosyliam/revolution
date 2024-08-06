package common

type Key int
type Receiver <-chan struct{}

var Send = struct{}{}

const (
	Forward Key = iota
	BackKey
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

type InterruptType int

const (
	ImmediateInterrupt InterruptType = iota
	DelayedInterrupt
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
	Execute(interruptType InterruptType) error
}

type Event interface {
	Execute(backend Backend)
}
