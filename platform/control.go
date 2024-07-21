//go:build darwin

package platform

// #include "control.h"
import "C"

import (
	"github.com/nosyliam/revolution/pkg/common"
	"unsafe"
)

var ControlBackend common.Backend = &controlBackend{}

type controlBackend struct{}

func (c controlBackend) KeyDown(pid int, key common.Key) {
	C.send_key_event(C.int(pid), true, C.int(KeyCodeMap[key]))
}

func (c controlBackend) KeyUp(pid int, key common.Key) {
	C.send_key_event(C.int(pid), false, C.int(KeyCodeMap[key]))
}

func (c controlBackend) MoveMouse(x, y int) {
	C.move_mouse(C.int(x), C.int(y))
}

func (c controlBackend) ScrollMouse(x, y int) {
	C.scroll_mouse(C.int(x), C.int(y))
}

func (c controlBackend) SleepAsync(ms int, interrupt <-chan struct{}) <-chan struct{} {
	sleepDone := make(chan struct{})

	go func() {
		intV := 0
		done := make(chan bool)

		go func() {
			select {
			case <-interrupt:
				intV = 1
			case <-done:
			}
		}()

		C.microsleep(C.int(ms), (*C.int)(unsafe.Pointer(&intV)))
		close(done)
		sleepDone <- struct{}{}
	}()

	return sleepDone
}

func (c controlBackend) Sleep(ms int, interrupt <-chan struct{}) {
	intV := 0
	done := make(chan bool)

	go func() {
		select {
		case <-interrupt:
			intV = 1
		case <-done:
		}
	}()

	C.microsleep(C.int(ms), (*C.int)(unsafe.Pointer(&intV)))
	close(done)
}
