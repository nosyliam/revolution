package main

import (
	"embed"
	"github.com/nosyliam/revolution/pkg/control"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/nosyliam/revolution/platform"
	"github.com/pkg/errors"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	windowManager := window.NewWindowManager(platform.WindowBackend)
	eventBus := control.NewEventBus(platform.ControlBackend)
	app := NewMacro(windowManager, eventBus)

	if err := wails.Run(&options.App{
		Title:         "Revolution Macro",
		Width:         600,
		Height:        400,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	}); err != nil {
		dialog.Message(errors.Wrap(err, "Failed to run application").Error()).Error()
		return
	}
}
