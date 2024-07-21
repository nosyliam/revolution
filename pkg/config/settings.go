package config

type DiscordSettings struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookUrl string `yaml:"webhookUrl"`
	PingID     int64  `yaml:"pingID"`
}

// Settings defines the configuration for an individual preset
type Settings struct {
	Name         string           `yaml:"name"`
	LogVerbosity int              `yaml:"logVerbosity"`
	Discord      *DiscordSettings `yaml:"discord"`

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
	Presets []*Settings `yaml:"presets"`
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
