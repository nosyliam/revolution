//go:build windows
// +build windows

package networking

import (
	"github.com/sqweek/dialog"
	"syscall"
)

var ComputerName string
var UserName string

func getIdentity() string {
	return ComputerName + "/" + UserName
}

func getComputerName() (string, error) {
	var size uint32 = 256
	buf := make([]uint16, size)
	err := syscall.GetComputerName(&buf[0], &size)
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(buf[:size]), nil
}

func getUserName() (string, error) {
	var size uint32 = 256
	buf := make([]uint16, size)
	err := syscall.GetUserName(&buf[0], &size)
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(buf[:size]), nil
}

func init() {
	var err error
	ComputerName, err = getMacComputerName()
	if err != nil {
		dialog.Message("Failed to get computer name: %s", err.Error()).Error()
		os.Exit(1)
	}
	UserName, err = getMacUserName()
	if err != nil {
		dialog.Message("Failed to get user name: %s", err.Error()).Error()
		os.Exit(1)
	}
}
