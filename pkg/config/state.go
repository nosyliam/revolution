package config

type MacroState struct {
	state *State
}

func (m *MacroState) Save() error {
	return m.state.Save()
}

type State struct {
	configFile
	Macros map[string]*MacroState `json:"macros"`
}

func NewState() *State {
	return &State{configFile: configFile{path: "state.yaml", format: YAML}, Macros: make(map[string]*MacroState)}
}
