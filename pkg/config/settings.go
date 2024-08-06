package config

import "github.com/pkg/errors"

type WindowAlignment string

const (
	TopLeftWindowAlignment     WindowAlignment = "top-left"
	TopRightWindowAlignment                    = "top-right"
	BottomLeftWindowAlignment                  = "bottom-left"
	BottomRightWindowAlignment                 = "bottom-right"
)

type DiscordSettings struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookUrl string `yaml:"webhookUrl"`
	PingID     int64  `yaml:"pingID"`
}

type WindowConfig struct {
	ID        string          `yaml:"id"`
	Alignment WindowAlignment `yaml:"alignment"`
	FullWidth bool            `yaml:"fullWidth"`
	Screen    int             `yaml:"screen"`
}

// Settings defines the configuration for an individual preset
type Settings struct {
	Name         string           `yaml:"name"`
	AccountName  *string          `yaml:"accountName"`
	LogVerbosity int              `yaml:"logVerbosity"`
	Discord      *DiscordSettings `yaml:"discord"`

	WindowConfigID         *string `yaml:"windowConfigId"`
	PrivateServerLink      *string `yaml:"privateServerLink"`
	FallbackToPublicServer bool    `yaml:"fallbackToPublicServer"`

	config *Config
}

func (s *Settings) WindowConfig(id string) *WindowConfig {
	for _, window := range s.config.Windows {
		if window.ID == id {
			return window
		}
	}
	return nil
}

func (s *Settings) Default() {
	s.Discord = &DiscordSettings{}
	s.FallbackToPublicServer = true
}

func (s *Settings) Save() error {
	return s.config.Save()
}

type Config struct {
	configFile
	Presets []*Settings     `yaml:"presets"`
	Windows []*WindowConfig `yaml:"windows"`
}

func (c *Config) NewPreset(name string) *Settings {
	preset := &Settings{config: c}
	preset.Default()
	c.Presets = append(c.Presets, preset)
	return preset
}

func NewConfig() (*Config, error) {
	config := &Config{configFile: configFile{path: "settings.yaml", format: YAML}}
	if err := config.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro config")
	}
	return config, nil
}
