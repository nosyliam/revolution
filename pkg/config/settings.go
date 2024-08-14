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
	configObject
	Enabled    bool   `yaml:"enabled"`
	WebhookUrl string `yaml:"webhookUrl"`
	PingID     int64  `yaml:"pingID"`
}

type WindowConfig struct {
	configObject
	ID        string          `yaml:"id"`
	Alignment WindowAlignment `yaml:"alignment"`
	FullWidth bool            `yaml:"fullWidth"`
	Screen    int             `yaml:"screen"`
}

func (w *WindowConfig) Key() string {
	return w.ID
}

type WindowSettings struct {
	WindowConfigID         *string     `yaml:"windowConfigId"`
	DefaultWindowSize      *WindowSize `yaml:"windowSize"`
	PrivateServerLink      *string     `yaml:"privateServerLink"`
	FallbackToPublicServer bool        `yaml:"fallbackToPublicServer"`
}

func (c *WindowSettings) Default() {
	c.FallbackToPublicServer = false
}

// Settings defines the configuration for an individual preset
type Settings struct {
	configObject
	Name         string           `yaml:"name"`
	AccountName  *string          `yaml:"accountName"`
	LogVerbosity int              `yaml:"logVerbosity"`
	Discord      *DiscordSettings `yaml:"discord"`
	Window       *WindowSettings  `yaml:"window"`

	config *Config
}

func (c *Settings) Key() string {
	return c.Name
}

func (s *Settings) WindowConfig(id string) *WindowConfig {
	for _, window := range s.config.Windows.data {
		if window.ID == id {
			return window
		}
	}
	return nil
}

func (s *Settings) Save() error {
	return s.config.Save()
}

type Config struct {
	configFile
	configObject
	Presets *configList[*Settings]     `yaml:"presets"`
	Windows *configList[*WindowConfig] `yaml:"windows"`
}

func (c *Config) NewPreset(name string) *Settings {
	preset := &Settings{Name: name, config: c}
	preset.Initialize(name)
	c.Presets = append(c.Presets, preset)
	return preset
}

func NewConfig() (*Config, error) {
	config := &Config{configFile: configFile{path: "settings.yaml", format: YAML}}
	if err := config.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro config")
	}
	if len(config.Presets) == 0 {
		config.NewPreset("default")
	}
	return config, nil
}
