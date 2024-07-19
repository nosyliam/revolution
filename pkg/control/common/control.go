package common

import (
	"github.com/nosyliam/revolution/pkg/config"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/nosyliam/revolution/pkg/window"
)

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

type LogicResult int

const (
	ContinueLogicResult LogicResult = iota
	RetryLogicResult
	ErrorLogicResult
)

type ActionResults struct {
	imageSearchPoints []revimg.Point
}

type Dependencies struct {
	EventBus EventBus
	Backend  Backend
	Settings *config.Settings
	Logger   *logging.Logger
	Window   *window.Window
	Results  *ActionResults
	Exec     func(string) error
}

type Backend interface {
	KeyDown(pid int, key Key)
	KeyUp(pid int, key Key)
	MoveMouse(x, y int)
	ScrollMouse(x, y int)
	Sleep(ms int, interrupt <-chan struct{})
	SleepAsync(ms int, interrupt <-chan struct{}) <-chan struct{}
}

type EventBus interface {
	Start()
	KeyDown(pid int, key Key) <-chan struct{}
	KeyUp(pid int, key Key) <-chan struct{}
	MoveMouse(x, y int) <-chan struct{}
	ScrollMouse(x, y int) <-chan struct{}
}

type Action interface {
	Execute(deps *Dependencies) error
}

type Event interface {
	Execute(id int, backend Backend)
}
