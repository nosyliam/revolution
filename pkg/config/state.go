package config

type MacroState struct {
	state *State

	Unwind bool
}

func (m *MacroState) Save() error {
	return m.state.Save()
}

type State struct {
	configFile
	Macros map[string]*MacroState `json:"macros"`
}

func (s *State) State(name string) *MacroState {
	if state, ok := s.Macros[name]; ok {
		return state
	}
	state := &MacroState{state: s}
	s.Macros[name] = state
	return state
}

func NewState() *State {
	return &State{configFile: configFile{path: "state.yaml", format: YAML}, Macros: make(map[string]*MacroState)}
}
