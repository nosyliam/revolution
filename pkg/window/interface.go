package window

import (
	revimage "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
)

var (
	PermissionDeniedErr = errors.New("accessibility permissions required")
	WindowNotFoundErr   = errors.New("window not found")
)
var Manager = &windowManager{}

type JoinOptions struct {
	LinkCode string
	Url      string
}

type Backend interface {
	OpenWindow(options JoinOptions) (int, error)
	CloseWindow(id int) error
	ActivateWindow(id int) error
	SetRobloxLocation(loc string)

	Screenshot(id int) (*image.RGBA, error)
	GetFrame(id int) (*revimage.Frame, error)
	SetFrame(id int, frame revimage.Frame) error
	DisplayFrames() ([]revimage.Frame, error)
}

type Window struct {
	backend Backend
	id      int
}

func (w *Window) FindImage(bitmapName string) int {
	return 0
}

type windowManager struct {
	backend Backend
}

func (m *windowManager) FindRoblox() (*Window, error) {
	return nil, nil
}

//func (m *windowManager) OpenRoblox() *Window
