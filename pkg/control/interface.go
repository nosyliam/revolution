package control

type Key int

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

type Backend interface {
	KeyDown(pid int, key Key)
	KeyUp(pid int, key Key)
	MoveMouse(x, y int)
	ScrollMouse(x, y int)
	Sleep(ms int, interrupt <-chan interface{})
	SleepAsync(ms int, interrupt <-chan interface{}) <-chan interface{}
}
