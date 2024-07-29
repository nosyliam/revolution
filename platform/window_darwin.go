//go:build darwin

package platform

// #include "window.h"
import "C"

import "C"

import (
	"encoding/json"
	"fmt"
	ps "github.com/mitchellh/go-ps"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/pkg/errors"
	"image"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"
	"unsafe"
)

var WindowBackend window.Backend = &windowBackend{make(map[int]*C.Window), ""}

// This structure is redundant because it's impossible to perform multi-instance on MacOS, but that may change
type windowBackend struct {
	windows   map[int]*C.Window
	robloxLoc string
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
	return nil
}

func (w *windowBackend) freeWindow(id int) {
	win, ok := w.windows[id]
	if !ok {
		return
	}

	C.CFRelease((C.CFTypeRef)(win.window))
	C.free(unsafe.Pointer(win))
}

func (w *windowBackend) getWindow(id int) (*C.Window, error) {
	resp := (bool)(C.check_ax_enabled(C.bool(true)))
	if resp == false {
		return nil, window.PermissionDeniedErr
	}

	win, ok := w.windows[id]
	if !ok {
		return nil, window.WindowNotFoundErr
	}

	_, err := ps.FindProcess(id)
	if err != nil {
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

func (w *windowBackend) getRobloxVersion(loc string) (string, error) {
	pList, err := os.Open(fmt.Sprintf("%s/Contents/Info.plist", loc))
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "failed to open Roblox application info for update check")
	}

	data, err := io.ReadAll(pList)
	if err != nil {
		return "", errors.Wrap(err, "failed to read Roblox application info for update check")
	}

	_ = pList.Close()

	versionReg := regexp.MustCompile("<key>CFBundleShortVersionString</key>\n\t<string>(.*)</string>")
	matches := versionReg.FindStringSubmatch(string(data))
	if len(matches) <= 1 {
		return "", errors.Wrap(err, "failed to read Roblox version")
	}

	return matches[1], nil
}

func (w *windowBackend) SetRobloxLocation(loc string) {
	w.robloxLoc = loc
}

func (w *windowBackend) OpenWindow(options window.JoinOptions) (int, error) {
	resp := (bool)(C.check_ax_enabled(C.bool(true)))
	if resp == false {
		return 0, window.PermissionDeniedErr
	}

	// Attempt to bind non-macro roblox instances first
	findRunningInstance := func(initialize bool) (int, error) {
		processes, err := ps.Processes()
		if err != nil {
			return 0, errors.Wrap(err, "failed to list processes")
		}
		for _, process := range processes {
			if pid := process.Pid(); process.Executable() == "RobloxPlayer" {
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

	loc := "/Applications/Roblox.app"
	if w.robloxLoc != "" {
		loc = w.robloxLoc
	}

	version, err := w.getRobloxVersion(loc)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get current Roblox version")
	}

	vResp, err := http.Get("https://clientsettingscdn.roblox.com/v2/client-version/MacPlayer")
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

	if dataMap["version"].(string) != version {
		var pid = -1
		cmd := exec.Command("open", fmt.Sprintf("%s/Contents/MacOS/RobloxPlayerInstaller.app", loc))
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
			if int((C.int)(C.get_window_count(C.int(pid)))) == 1 {
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
		cmd := exec.Command("open", "-n", url)
		err = cmd.Start()
		if err != nil {
			return 0, errors.New("failed to start installer")
		}

		for j := 0; j < 100; j++ {
			if pid, err := findRunningInstance(true); err != nil {
				return 0, err
			} else if pid != -1 {
				return pid, nil
			}
			<-time.After(10 * time.Millisecond)
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
	data := C.GoBytes(unsafe.Pointer(screen.data), C.int(screen.len))
	width, height, stride := C.ulong(screen.width), C.ulong(screen.height), C.ulong(screen.stride)

	img := image.RGBA{}
	img.Rect = image.Rect(0, 0, int(width), int(height))
	img.Pix = data
	img.Stride = int(stride)

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
	if cFrame == nil {
		return nil, errors.New("failed to get window frame")
	}

	return &revimg.Frame{
		Width:  int(C.int(cFrame.width)),
		Height: int(C.int(cFrame.height)),
		X:      int(C.int(cFrame.x)),
		Y:      int(C.int(cFrame.y)),
	}, nil
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
