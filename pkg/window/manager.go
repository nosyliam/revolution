package window

import (
	"fmt"
	"github.com/nosyliam/revolution/bitmaps"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
	"sync"
)

type Window struct {
	id         int
	screen     int
	settings   *config.WindowSettings
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
	if errors.Is(err, WindowNotFoundErr) {
		w.mgr.freeWindow(w.settings.ID)
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func (w *Window) Screenshot() *image.RGBA {
	return w.screenshot
}

func (w *Window) TakeScreenshot() error {
	var err error
	w.screenshot, err = w.backend.Screenshot(w.id)
	if err != nil {
		return err
	}
	return nil
}

func (w *Window) Close() error {
	return w.backend.CloseWindow(w.id)
}

type Manager struct {
	sync.Mutex
	backend         Backend
	reservedWindows [][4]*string
	reservedIds     map[string]bool
	windowFrames    map[string]revimg.Frame
	frames          []revimg.ScreenFrame
}

func (m *Manager) freeWindow(id string) {

}

func (m *Manager) windowFrame(id string) revimg.Frame {
	return m.windowFrames[id]
}

func (m *Manager) reserveWindow(settings config.WindowSettings) error {
	m.Lock()
	defer m.Unlock()
	screens, err := m.backend.DisplayFrames()
	if err != nil {
		return errors.Wrap(err, "failed to get screen data")
	}
	if len(screens) < settings.Screen {
		return errors.New(fmt.Sprintf("The configuration for this window exists on screen %d. Only %d screen(s) are available.",
			settings.Screen, len(screens)))
	}
	m.frames = screens

	screen := screens[settings.Screen]
	if screen.Scale > 1 {
		return errors.New(fmt.Sprintf("Retina displays are not supported. Attach a monitor or install DeskPad to add a virtual monitor."))
	}

	if len(screens) > len(m.reservedWindows) {
		for i := len(m.reservedWindows); i < len(screens); i++ {
			m.reservedWindows = append(m.reservedWindows, [4]*string{})
		}
	}
	reservations := m.reservedWindows[settings.Screen]
	id := settings.ID

	// Validate and reserve the position/size
	switch settings.Alignment {
	case config.TopLeftWindowAlignment:
		if reservations[0] != nil {
			return errors.New(fmt.Sprintf("The top left corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[1] != nil {
			return errors.New(fmt.Sprintf("The top right corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[1] = &id
		}
		reservations[0] = &id
	case config.TopRightWindowAlignment:
		if reservations[1] != nil {
			return errors.New(fmt.Sprintf("The top right corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[0] != nil {
			return errors.New(fmt.Sprintf("The top left corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[0] = &id
		}
		reservations[1] = &id
	case config.BottomLeftWindowAlignment:
		if reservations[2] != nil {
			return errors.New(fmt.Sprintf("The bottom left corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[3] != nil {
			return errors.New(fmt.Sprintf("The bottom right corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[3] = &id
		}
		reservations[2] = &id
	case config.BottomRightWindowAlignment:
		if reservations[3] != nil {
			return errors.New(fmt.Sprintf("The bottom right corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[2] != nil {
			return errors.New(fmt.Sprintf("The bottom left corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[2] = &id
		}
		reservations[3] = &id
	}

	// Compute the frame of the window
	x := 0
	for i := 0; i < settings.Screen; i++ {
		x += screens[i].Width
	}
	y := 0
	w := screen.Width
	if !settings.FullWidth {
		w = screen.Width / 2
	}
	if settings.Alignment == config.BottomLeftWindowAlignment || settings.Alignment == config.BottomRightWindowAlignment {
		y = screen.Height / 2
	}
	if (settings.Alignment == config.TopRightWindowAlignment || settings.Alignment == config.BottomRightWindowAlignment) && settings.FullWidth == false {
		x += screen.Width / 2
	}

	m.windowFrames[settings.ID] = revimg.Frame{X: x, Y: y, Width: w, Height: screen.Height / 2}
	return nil
}

func (m *Manager) OpenWindow(macro *common.Macro) (*Window, error) {
	return nil, nil
}

func NewWindowManager(backend Backend) *Manager {
	return &Manager{backend: backend}
}
