package window

import (
	"github.com/nosyliam/revolution/bitmaps"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
)

var (
	PermissionDeniedErr = errors.New("accessibility permissions required")
	WindowNotFoundErr   = errors.New("window not found")
)

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
	GetFrame(id int) (*revimg.Frame, error)
	SetFrame(id int, frame revimg.Frame) error
	DisplayFrames() ([]revimg.Frame, error)
}

type Window struct {
	backend    Backend
	id         int
	screenshot *image.RGBA
}

func (w *Window) PID() int {
	return w.id
}

func (w *Window) FindImage(bitmapName string, options *revimg.SearchOptions) ([]revimg.Point, error) {
	needle := bitmaps.Registry.Get(bitmapName)
	return revimg.ImageSearch(needle, w.screenshot, options)
}

func (w *Window) Screenshot() error {
	var err error
	w.screenshot, err = w.backend.Screenshot(w.id)
	if err != nil {
		return err
	}
	return nil
}

type Manager struct {
	backend Backend
}

func (m *Manager) OpenRoblox() (*Window, error) {
	return nil, nil
}

func NewWindowManager(backend Backend) *Manager {
	return &Manager{backend: backend}
}
