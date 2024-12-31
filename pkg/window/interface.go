package window

import (
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
	DissociateWindow(int)
	OpenWindow(options JoinOptions) (int, error)
	CloseWindow(id int) error
	ActivateWindow(id int) error
	SetRobloxLocation(loc string)

	StartCapture(id int) (<-chan *image.RGBA, error)
	StopCapture(id int)
	GetFrame(id int) (*revimg.Frame, error)
	SetFrame(id int, frame revimg.Frame) error
	DisplayFrames() ([]revimg.ScreenFrame, error)
	DisplayCount() (int, error)
}
