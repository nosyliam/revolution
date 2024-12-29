//go:build windows
// +build windows

package platform

// #include "window.h"
// #cgo LDFLAGS: -lgdi32 -lshcore
import "C"

import (
	"encoding/json"
	"fmt"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
	ps "github.com/shirou/gopsutil/v4/process"
	"image"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"
)

var WindowBackend window.Backend = &windowBackend{make(map[int]*C.Window), ""}

type windowBackend struct {
	windows   map[int]*C.Window
	robloxLoc string
}

func (w *windowBackend) DissociateWindow(id int) {
	if _, ok := w.windows[id]; ok {
		w.freeWindow(id)
	}
	delete(w.windows, id)
}

func (w *windowBackend) CloseWindow(id int) error {
	proc, err := os.FindProcess(id)
	if err == nil {
		err = proc.Kill()
		if err != nil {
			return err
		}
	}

	w.freeWindow(id)
	delete(w.windows, id)
	return nil
}

func (w *windowBackend) freeWindow(id int) {
	win, ok := w.windows[id]
	if !ok {
		return
	}

	C.free(unsafe.Pointer(win))
}

func (w *windowBackend) getWindow(id int) (*C.Window, error) {
	win, ok := w.windows[id]
	if !ok {
		return nil, window.WindowNotFoundErr
	}

	exists, err := ps.PidExists(int32(id))
	if !exists || err != nil {
		w.freeWindow(id)
		return nil, window.WindowNotFoundErr
	}

	return win, nil
}

func (w *windowBackend) initializeWindow(pid int) bool {
	ret := (*C.Window)(C.get_window_with_pid(C.int(pid)))
	if ret == nil {
		return false
	}

	w.windows[pid] = ret
	return true
}

func (w *windowBackend) getRobloxRoot() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "Roblox", "Versions")
}

func (w *windowBackend) getRobloxVersion() (string, error) {
	entries, err := os.ReadDir(w.getRobloxRoot())
	if err != nil {
		panic(err)
	}

	var instances = make(map[string]time.Time)
	for _, entry := range entries {
		if entry.IsDir() {
			exePath := filepath.Join(w.getRobloxRoot(), entry.Name(), "RobloxPlayerBeta.exe")
			if info, err := os.Stat(exePath); err == nil {
				instances[entry.Name()] = info.ModTime()
			} else if !os.IsNotExist(err) {
				return "", err
			}
		}
	}

	if len(instances) == 0 {
		return "", errors.New("no Roblox versions found")
	}

	var mostRecent time.Time
	var mostRecentKey string

	for key, t := range instances {
		if mostRecent.IsZero() || t.After(mostRecent) {
			mostRecent = t
			mostRecentKey = key
		}
	}

	return mostRecentKey, nil
}

func (w *windowBackend) SetRobloxLocation(loc string) {
	w.robloxLoc = loc
}

func (w *windowBackend) OpenWindow(options window.JoinOptions) (int, error) {
	// Attempt to bind non-macro roblox instances first
	findRunningInstance := func(initialize bool) (int, error) {
		processes, err := ps.Processes()
		if err != nil {
			return 0, errors.Wrap(err, "failed to list processes")
		}
		for _, p := range processes {
			n, err := p.Name()
			if err != nil {
				continue
			}
			if pid := int(p.Pid); n == "RobloxPlayerBeta.exe" {
				if !initialize {
					return pid, nil
				}
				if _, ok := w.windows[pid]; !ok {
					if w.initializeWindow(pid) {
						return pid, nil
					}
				}
			}
		}
		return -1, nil
	}

	if pid, err := findRunningInstance(true); err != nil {
		return 0, err
	} else if pid != -1 {
		return pid, nil
	}

	version, err := w.getRobloxVersion()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get current Roblox version")
	}

	vResp, err := http.Get("https://clientsettingscdn.roblox.com/v2/client-version/WindowsPlayer")
	if err != nil {
		return 0, errors.Wrap(err, "failed to fetch latest Roblox version")
	}

	data, _ := io.ReadAll(vResp.Body)
	fmt.Println(string(data))

	dataMap := make(map[string]interface{})
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal latest Roblox version")
	}

	if dataMap["clientVersionUpload"].(string) != version {
		var pid = -1
		cmd := exec.Command(filepath.Join(w.getRobloxRoot(), "RobloxPlayerInstaller.exe"))
		err = cmd.Start()
		if err != nil {
			return 0, errors.New("failed to start installer")
		}

		for i := 0; i < 1000; i++ {
			if pid, err = findRunningInstance(false); err != nil {
				return 0, err
			} else if pid != -1 {
				break
			}
			<-time.After(10 * time.Millisecond)
		}

		if pid == -1 {
			return 0, errors.New("failed to find installer process")
		}

		for i := 0; i < 1000; i++ {
			if int((C.int)(C.get_window_visible_count(C.int(pid)))) == 1 {
				break
			}
			<-time.After(10 * time.Millisecond)
		}

		proc, err := os.FindProcess(pid)
		if err == nil {
			_ = proc.Kill()
		}
	}

	var url string
	if options.Url == "" {
		url = fmt.Sprintf("roblox://placeID=1537690962%s", (func() string {
			if options.LinkCode == "" {
				return ""
			} else {
				return fmt.Sprintf("&linkCode=%s", options.LinkCode)
			}
		})())
	} else {
		url = options.Url
	}

	for i := 0; i < 10; i++ {
		cmd := exec.Command("cmd", "/C", "start", url)
		err = cmd.Start()
		fmt.Println("starting roblox")
		if err != nil {
			return 0, errors.New("failed to open roblox")
		}

		var pid = 0
		for j := 0; j < 500; j++ {
			if pid, err = findRunningInstance(true); err != nil {
				return 0, err
			} else if pid != -1 {
				return pid, nil
			}
			<-time.After(10 * time.Millisecond)
		}

		if pid == -1 {
			return 0, errors.New("failed to find roblox process")
		}

		<-time.After(1 * time.Second)
	}

	return 0, errors.New("failed to initialize roblox process")
}

func (w *windowBackend) ActivateWindow(id int) error {
	win, err := w.getWindow(id)
	if err != nil {
		return err
	}

	C.activate_window(win)
	return nil
}

func (w *windowBackend) Screenshot(id int) (*image.RGBA, error) {
	win, err := w.getWindow(id)
	if err != nil {
		return nil, err
	}

	screen := (*C.Screenshot)(C.screenshot_window(win))
	if screen == nil {
		return nil, errors.New("failed to take screenshot")
	}
	data := C.GoBytes(unsafe.Pointer(screen.data), C.int(screen.len))
	width, height, stride := C.ulong(screen.width), C.ulong(screen.height), C.ulong(screen.stride)

	img := image.RGBA{}
	img.Rect = image.Rect(0, 0, int(width), int(height))
	img.Pix = data
	img.Stride = int(stride)

	/*f, _ := os.Create("test.png")
	png.Encode(f, &img)
	f.Close()*/

	C.free(unsafe.Pointer(screen.data))
	C.free(unsafe.Pointer(screen))
	return &img, nil
}

func (w *windowBackend) GetFrame(id int) (*revimg.Frame, error) {
	win, err := w.getWindow(id)
	if err != nil {
		return nil, err
	}

	cFrame := (*C.Frame)(C.get_window_frame(win))
	if cFrame == nil || (cFrame.width == 0 && cFrame.height == 0) {
		return nil, errors.New("failed to get window frame")
	}

	frame := &revimg.Frame{
		Width:  int(C.int(cFrame.width)),
		Height: int(C.int(cFrame.height)),
		X:      int(C.int(cFrame.x)),
		Y:      int(C.int(cFrame.y)),
	}
	C.free(unsafe.Pointer(cFrame))
	return frame, nil
}

func (w *windowBackend) SetFrame(id int, frame revimg.Frame) error {
	win, err := w.getWindow(id)
	if err != nil {
		return err
	}

	C.set_window_frame(win, C.int(frame.Width), C.int(frame.Height), C.int(frame.X), C.int(frame.Y))
	return nil
}

func (w *windowBackend) DisplayFrames() ([]revimg.ScreenFrame, error) {
	var goFrames []revimg.ScreenFrame
	frames := (*C.Frames)(C.get_display_frames())
	if int(C.int(frames.len)) == 0 {
		return nil, errors.New("no displays were detected")
	}
	for i := 0; i < int(C.int(frames.len)); i++ {
		cFrame := (*C.Frame)(unsafe.Pointer(uintptr(unsafe.Pointer(frames.frames)) + (unsafe.Sizeof(C.Frame{}) * uintptr(i))))
		goFrames = append(goFrames, revimg.ScreenFrame{
			Frame: revimg.Frame{
				Width:  int(C.int(cFrame.width)),
				Height: int(C.int(cFrame.height)),
				X:      int(C.int(cFrame.x)),
				Y:      int(C.int(cFrame.y)),
			},
			Scale: float32(C.float(cFrame.scale)),
		})
	}
	C.free(unsafe.Pointer(frames))
	return goFrames, nil
}

func (w *windowBackend) DisplayCount() (int, error) {
	res := (C.int)(C.get_display_count())
	if count := int(res); count == -1 {
		return 0, errors.New("failed to read display information")
	} else {
		return count, nil
	}
}
