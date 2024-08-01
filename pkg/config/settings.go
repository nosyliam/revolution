package config

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

type WindowSettings struct {
	ID        string          `yaml:"id"`
	Alignment WindowAlignment `yaml:"alignment"`
	FullWidth bool            `yaml:"fullWidth"`
	Screen    int             `yaml:"screen"`
}

// Settings defines the configuration for an individual preset
type Settings struct {
	Name           string           `yaml:"name"`
	WindowConfigID *string          `yaml:"windowId"`
	LogVerbosity   int              `yaml:"logVerbosity"`
	Discord        *DiscordSettings `yaml:"discord"`

	config *Config
}

func (s *Settings) Default() {
	s.Discord = &DiscordSettings{}
}

func (s *Settings) Save() error {
	return s.config.Save()
}

type Config struct {
	configFile
	Presets []*Settings       `yaml:"presets"`
	Windows []*WindowSettings `yaml:"windows"`
}

func (c *Config) NewPreset(name string) *Settings {
	preset := &Settings{config: c}
	preset.Default()
	c.Presets = append(c.Presets, preset)
	return preset
}

func NewConfig() *Config {
	return &Config{configFile: configFile{path: "settings.yaml", format: YAML}}
}
