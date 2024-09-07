package config

import (
	"github.com/pkg/errors"
)

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
	WebhookUrl string `yaml:"webhookUrl,omitempty"`
	PingID     int    `yaml:"pingID,omitempty"`
}

type WindowConfig struct {
	ID        string          `yaml:"id" key:"true" lock:"default"`
	Alignment WindowAlignment `yaml:"alignment" default:"top-left"`
	FullWidth bool            `yaml:"fullWidth" default:"true"`
	Screen    int             `yaml:"screen"`
}

type WindowSettings struct {
	WindowConfigID         string     `yaml:"windowConfigId"`
	DefaultWindowSize      WindowSize `yaml:"windowSize" default:"fullscreen"`
	PrivateServerLink      string     `yaml:"privateServerLink,omitempty"`
	FallbackToPublicServer bool       `yaml:"fallbackToPublicServer" default:"true"`
}

// Settings defines the configuration for an individual preset
type Settings struct {
	Name         string                   `yaml:"name" key:"true" lock:"default"`
	LogVerbosity int                      `yaml:"logVerbosity"`
	Discord      *Object[DiscordSettings] `yaml:"discord"`
	Windows      *Object[WindowSettings]  `yaml:"windows"`
}

type Tools struct {
	JellyTool *Object[JellyTool] `yaml:"jellyTool"`
}

type Config struct {
	Test    int                 `yaml:"test"`
	Presets *List[Settings]     `yaml:"presets"`
	Windows *List[WindowConfig] `yaml:"windows"`
	Tools   *Object[Tools]      `yaml:"tools"`
}

func NewConfig(runtime *Runtime) (Reactive, error) {
	settings := File[Config]{name: "settings", path: "settings.yaml", format: YAML, runtime: runtime}
	if err := settings.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro settings")
	}
	obj := settings.Object()
	if obj.LengthPath("presets") == 0 {
		_ = obj.AppendPath("presets[default]")
	}
	if obj.LengthPath("windows") == 0 {
		_ = obj.AppendPath("windows[default]")
	}
	return obj, nil
}
