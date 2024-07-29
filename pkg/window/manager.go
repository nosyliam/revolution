package window

import (
	"fmt"
	"github.com/nosyliam/revolution/bitmaps"
	"github.com/nosyliam/revolution/pkg/config"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
)

type Window struct {
	id         int
	screen     int
	frame      revimg.Frame
	backend    Backend
	screenshot *image.RGBA
	mgr        *Manager
}

func (w *Window) PID() int {
	return w.id
}

func (w *Window) FindImage(bitmapName string, options *revimg.SearchOptions) ([]revimg.Point, error) {
	needle := bitmaps.Registry.Get(bitmapName)
	return revimg.ImageSearch(needle, w.screenshot, options)
}

func (w *Window) Fix() error {
	_, err := w.backend.GetFrame(w.id)
	if err != nil {
		return err
	}
	return nil
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
	backend         Backend
	reservedWindows [][4]*string
	reservedIds     map[string]bool
}

func (m *Manager) ReserveWindow(settings config.WindowSettings) error {
	screens, err := m.backend.DisplayFrames()
	if err != nil {
		return errors.Wrap(err, "failed to get screen data")
	}
	if len(screens) < settings.Screen {
		return errors.New(fmt.Sprintf("The configuration for this window exists on screen %d. Only %d screen(s) are available.",
			settings.Screen, len(screens)))
	}
	return nil
}

func (m *Manager) OpenWindow() (*Window, error) {
	return nil, nil
}

func NewWindowManager(backend Backend) *Manager {
	return &Manager{backend: backend}
}
