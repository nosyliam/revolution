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
)

type System interface {
	SendKeyDown(Key) error
	SendKeyUp(Key) error
	MoveMouse()
}
