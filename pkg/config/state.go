package config

import (
	"github.com/pkg/errors"
)

type UnwindLoop struct {
	Depth    int
	Continue bool
}

type LoopState struct {
	Unwind *UnwindLoop
	Index  []int
}

type MacroNetworkingConfig struct {
	AvailableRelays     *List[string] `state:"availableRelays" yaml:"-"`
	ConnectedIdentities *List[string] `state:"connectedIdentities" yaml:"-"`
	ConnectingIdentity  string        `state:"connectingIdentity" yaml:"-"`
	RelayIdentity       string        `state:"relayIdentity" yaml:"-"`
}

type MacroState struct {
	AccountName string `yaml:"accountName" key:"true"`

	Running bool `state:"running" yaml:"-"`
	Paused  bool `state:"paused" yaml:"-"`

	Status string `state:"status" default:"Ready" yaml:"-"`

	HoneyOriginX int `state:"honeyOriginX" yaml:"-"`
	BaseOriginX  int `state:"baseOriginX" yaml:"-"`
	BaseOriginY  int `state:"baseOriginY" yaml:"-"`

	Networking *MacroNetworkingConfig `yaml:"networking"`
}

type StateConfig struct {
	CodeStatus    string `state:"codeStatus" default:"pending" yaml:"-"`
	DefaultPreset string `yaml:"defaultPreset" default:"default"`
	ActiveAccount string `yaml:"activeAccount,omitempty" default:"default"`
}

type NetworkingConfig struct {
	HistoricalIdentities *List[string] `yaml:"historicalIdentities"`
}

type State struct {
	Config     *Object[StateConfig] `yaml:"config"`
	Macros     *List[MacroState]    `yaml:"macros"`
	Networking *NetworkingConfig    `yaml:"networking"`
}

func NewState(runtime *Runtime) (*Object[State], error) {
	state := File[State]{name: "state", path: "state.yaml", format: YAML, runtime: runtime}
	if err := state.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro state")
	}
	return state.Object(), nil
}
