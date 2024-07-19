package config

// Settings defines the configuration for an individual preset
type Settings struct {
	Name string

	config *Config
}

type Config struct {
	Presets []Settings `json:"presets"`
}
