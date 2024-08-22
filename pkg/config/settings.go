package config

import "github.com/pkg/errors"

type WindowAlignment string

const (
	FullScreenWindowAlignment  WindowAlignment = "fullscreen"
	TopLeftWindowAlignment                     = "top-left"
	TopRightWindowAlignment                    = "top-right"
	BottomLeftWindowAlignment                  = "bottom-left"
	BottomRightWindowAlignment                 = "bottom-right"
)

type WindowSize string

const (
	QuarterWindowSize WindowSize = "quarter"
	HalfWindowSize    WindowSize = "half"
	FullWindowSize    WindowSize = "full"
)

type DiscordSettings struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookUrl string `yaml:"webhookUrl"`
	PingID     int64  `yaml:"pingID"`
}

type WindowConfig struct {
	ID        string          `yaml:"id" key:"true"`
	Alignment WindowAlignment `yaml:"alignment"`
	FullWidth bool            `yaml:"fullWidth"`
	Screen    int             `yaml:"screen"`
}

type WindowSettings struct {
	WindowConfigID         *string     `yaml:"windowConfigId"`
	DefaultWindowSize      *WindowSize `yaml:"windowSize"`
	PrivateServerLink      *string     `yaml:"privateServerLink"`
	FallbackToPublicServer bool        `yaml:"fallbackToPublicServer" default:"true"`
}

// Settings defines the configuration for an individual preset
type Settings struct {
	Name         string                   `yaml:"name" key:"true"`
	LogVerbosity int                      `yaml:"logVerbosity"`
	Discord      *Object[DiscordSettings] `yaml:"discord"`
	Window       *Object[WindowSettings]  `yaml:"window"`
}

type Config struct {
	Presets *List[Settings]     `yaml:"presets"`
	Windows *List[WindowConfig] `yaml:"windows"`
}

func NewConfig() (*Object[Config], error) {
	settings := File[Config]{path: "settings.yaml", format: YAML}
	if err := settings.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro settings")
	}
	obj := settings.Object()
	if len(obj.Concrete().Presets.Concrete()) == 0 {
		_ = obj.AppendPath("presets")
	}
	return obj, nil
}
