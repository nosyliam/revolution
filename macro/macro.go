package macro

import (
	"github.com/nosyliam/revolution/pkg/control"
	"github.com/nosyliam/revolution/pkg/window"
)

type Macro struct {
	window *window.Window
	event  *control.EventBus
}
