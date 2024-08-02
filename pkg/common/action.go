package common

import (
	revimg "github.com/nosyliam/revolution/pkg/image"
	"image"
)

type Actions []Action

type ActionResults struct {
	EditedScreenshot  *image.RGBA
	ImageSearchPoints []revimg.Point
	RetryCount        int
}

type Action interface {
	Execute(macro *Macro) error
}
