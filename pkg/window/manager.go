package window

import (
	"fmt"
	"github.com/nosyliam/revolution/bitmaps"
	"github.com/nosyliam/revolution/pkg/config"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
	"strings"
	"sync"
)

type Window struct {
	id         int
	config     *config.WindowConfig
	backend    Backend
	screenshot *image.RGBA
	mgr        *Manager
	err        error
}

func (w *Window) PID() int {
	return w.id
}

func (w *Window) FindImage(bitmapName string, options *revimg.SearchOptions) ([]revimg.Point, error) {
	needle := bitmaps.Registry.Get(bitmapName)
	return revimg.ImageSearch(needle, w.screenshot, options)
}

func (w *Window) Fix() error {
	if err := w.mgr.adjustDisplays(); err != nil {
		return errors.Wrap(err, "failed to adjust displays")
	}
	if w.err != nil {
		return w.err
	}

	frame, err := w.backend.GetFrame(w.id)
	if errors.Is(err, WindowNotFoundErr) {
		w.mgr.freeWindow(w.config.ID)
		return err
	} else if err != nil {
		return errors.Wrap(err, "failed to get window frame")
	}

	if correct := w.mgr.windowFrames[w.config.ID]; !frame.Equals(correct) {
		if err = w.backend.SetFrame(w.id, correct); err != nil {
			return errors.Wrap(err, "failed to set window frame")
		}
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
	if err := w.backend.CloseWindow(w.id); err != nil {
		return errors.Wrap(err, "failed to close window")
	}
	w.mgr.freeWindow(w.config.ID)
	return nil
}

type windowArray [4]*string

func (a windowArray) Available(spots ...int) bool {
	return false
}

type Manager struct {
	sync.Mutex
	backend         Backend
	reservedWindows []windowArray
	reservedIds     map[string]*Window
	windowFrames    map[string]revimg.Frame
	displayCount    int
	frames          []revimg.ScreenFrame
}

func (m *Manager) freeWindow(id string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.reservedIds[id]; !ok {
		return
	}
	for _, screen := range m.reservedWindows {
		for i, win := range screen {
			if win != nil {
				if *win == id {
					screen[i] = nil
				}
			}
		}
	}
	delete(m.reservedIds, id)
	delete(m.windowFrames, id)
}

func (m *Manager) adjustDisplays() error {
	m.Lock()
	defer m.Unlock()
	if count, err := m.backend.DisplayCount(); err != nil {
		return errors.Wrap(err, "failed to read display count")
	} else {
		m.displayCount = count
	}
	if len(m.reservedWindows) > m.displayCount {
		for i := len(m.reservedWindows); i > m.displayCount; i-- {
			for _, id := range m.reservedWindows[i-1] {
				if id != nil {
					if window, ok := m.reservedIds[*id]; ok {
						if err := window.Close(); err != nil {
							return errors.Wrap(err, "failed to close window")
						}
						window.err = errors.New("The display for this window has been disconnected")
					}
				}
			}
		}
		m.reservedWindows = m.reservedWindows[:m.displayCount]
	} else if len(m.reservedWindows) < m.displayCount {
		for len(m.reservedWindows) != m.displayCount {
			m.reservedWindows = append(m.reservedWindows, [4]*string{})
		}
	}
	return nil
}

func (m *Manager) windowFrame(id string) revimg.Frame {
	return m.windowFrames[id]
}

func (m *Manager) reserveWindow(settings *config.WindowConfig, defaultSz *config.WindowSize) (*config.WindowConfig, error) {
	m.Lock()
	defer m.Unlock()

	screens, err := m.backend.DisplayFrames()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get screen data")
	}

	// If there's no window configuration associated with the macro instance, attempt to find an available spot
	if settings == nil {
		var sz = config.FullWindowSize
		if defaultSz != nil {
			sz = *defaultSz
		}
		for _, screen := range m.reservedWindows {
			switch true {
			case screen.Available(0, 1, 2, 3) && sz == config.FullWindowSize:
				settings = &config.WindowConfig{Alignment: config.FullScreenWindowAlignment}
			case screen.Available(0, 1) && sz == config.HalfWindowSize:
				settings = &config.WindowConfig{Alignment: config.TopLeftWindowAlignment, FullWidth: true}
			case screen.Available(2, 3) && sz == config.HalfWindowSize:
				settings = &config.WindowConfig{Alignment: config.BottomLeftWindowAlignment, FullWidth: true}
			case screen.Available(0):
				settings = &config.WindowConfig{Alignment: config.TopLeftWindowAlignment}
			case screen.Available(1):
				settings = &config.WindowConfig{Alignment: config.TopRightWindowAlignment}
			case screen.Available(2):
				settings = &config.WindowConfig{Alignment: config.BottomLeftWindowAlignment}
			case screen.Available(3):
				settings = &config.WindowConfig{Alignment: config.BottomRightWindowAlignment}
			}
			if settings != nil {
				break
			}
		}
		if settings != nil {
			return nil, errors.New("No displays available for the macro window")
		} else {

		}
	}

	if len(screens) < settings.Screen {
		return nil, errors.New(fmt.Sprintf("The configuration for this window exists on screen %d. Only %d screen(s) are available.",
			settings.Screen, len(screens)))
	}
	m.frames = screens

	screen := screens[settings.Screen]
	if screen.Scale > 1 {
		return nil, errors.New(fmt.Sprintf("Retina displays are not supported. Attach a monitor or install DeskPad to add a virtual monitor."))
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
	case config.FullScreenWindowAlignment:
		if reservations[0] != nil {
			return nil, errors.New(fmt.Sprintf("The top left corner of screen %d is reserved.", settings.Screen))
		}
		if reservations[1] != nil {
			return nil, errors.New(fmt.Sprintf("The top right corner of screen %d is reserved.", settings.Screen))
		}
		if reservations[2] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom left corner of screen %d is reserved.", settings.Screen))
		}
		if reservations[3] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom right corner of screen %d is reserved.", settings.Screen))
		}
	case config.TopLeftWindowAlignment:
		if reservations[0] != nil {
			return nil, errors.New(fmt.Sprintf("The top left corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[1] != nil {
			return nil, errors.New(fmt.Sprintf("The top right corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[1] = &id
		}
		reservations[0] = &id
	case config.TopRightWindowAlignment:
		if reservations[1] != nil {
			return nil, errors.New(fmt.Sprintf("The top right corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[0] != nil {
			return nil, errors.New(fmt.Sprintf("The top left corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[0] = &id
		}
		reservations[1] = &id
	case config.BottomLeftWindowAlignment:
		if reservations[2] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom left corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[3] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom right corner of screen %d is reserved.", settings.Screen))
		} else {
			reservations[3] = &id
		}
		reservations[2] = &id
	case config.BottomRightWindowAlignment:
		if reservations[3] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom right corner of screen %d is reserved.", settings.Screen))
		}
		if settings.FullWidth && reservations[2] != nil {
			return nil, errors.New(fmt.Sprintf("The bottom left corner of screen %d is reserved.", settings.Screen))
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
	return settings, nil
}

func (m *Manager) OpenWindow(accountName *string, settings *config.Settings, db *config.AccountDatabase, ignoreLink bool) (*Window, error) {
	var joinOptions = JoinOptions{}
	var windowConfig *config.WindowConfig
	if accountName != nil {
		account := db.Get(*accountName)
		if account == nil {
			return nil, errors.New(fmt.Sprintf("Account %s not found", *accountName))
		}
		if url, err := account.GenerateJoinUrl(ignoreLink); err != nil {
			return nil, errors.Wrap(err, "Failed to generate join url")
		} else {
			joinOptions = JoinOptions{Url: url}
		}
		if account.WindowConfigID != nil {
			if windowConfig = settings.WindowConfig(*account.WindowConfigID); windowConfig == nil {
				return nil, errors.New(fmt.Sprintf("Window configuration %s does not exist", *account.WindowConfigID))
			}
		}
	} else {
		if settings.Window.WindowConfigID != nil {
			if windowConfig = settings.WindowConfig(*settings.Window.WindowConfigID); windowConfig == nil {
				return nil, errors.New(fmt.Sprintf("Window configuration %s does not exist", *settings.Window.WindowConfigID))
			}
		}
		if !ignoreLink && settings.Window.PrivateServerLink != nil {
			parts := strings.Split(*settings.Window.PrivateServerLink, "=")
			if len(parts) != 2 {
				return nil, errors.New("Invalid private server link format")
			}
			joinOptions = JoinOptions{LinkCode: parts[1]}
		}
	}
	if windowConfig != nil {
		if _, ok := m.reservedIds[windowConfig.ID]; ok {
			return nil, errors.New(fmt.Sprintf("Window config ID \"%s\" is already in use", windowConfig.ID))
		}
	}
	var err error
	if err = m.adjustDisplays(); err != nil {
		return nil, errors.Wrap(err, "failed to adjust displays")
	}
	windowConfig, err = m.reserveWindow(windowConfig, settings.Window.DefaultWindowSize)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to reserve window")
	}
	id, err := m.backend.OpenWindow(joinOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open window")
	}

	win := &Window{
		id:      id,
		config:  windowConfig,
		backend: m.backend,
		mgr:     m,
	}
	m.reservedIds[windowConfig.ID] = win
	return win, nil
}

func NewWindowManager(backend Backend) *Manager {
	return &Manager{backend: backend}
}
