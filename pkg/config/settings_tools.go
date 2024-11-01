package config

type JellyTool struct {
	Enabled         bool          `yaml:"enabled"`
	BeeTypes        *List[string] `yaml:"beeTypes"`
	RequireMutation bool          `yaml:"requireMutation"`
	MutationType    string        `yaml:"mutationType" default:"Movespeed"`
	MutationValue   int           `yaml:"mutationValue" default:"0"`
	StopGifted      bool          `yaml:"stopGifted"`
	StopMythic      bool          `yaml:"stopMythic"`
}
