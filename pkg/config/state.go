package config

import "github.com/pkg/errors"

type UnwindLoop struct {
	Depth    int
	Continue bool
}

type LoopState struct {
	Unwind *UnwindLoop
	Index  []int
}

type MacroState struct {
	LoopState *LoopState
	LastError error
	Stack     []string

	PrivateServerAttempts int
	UsePublicServer       bool

	ClaimedHive *int

	state *State
}

func (m *MacroState) Save() error {
	return m.state.Save()
}

type State struct {
	configFile
	ActiveAccount *string                `json:"activeAccount"`
	Macros        map[string]*MacroState `json:"macros"`
}

func (s *State) State(name string) *MacroState {
	if state, ok := s.Macros[name]; ok {
		return state
	}
	state := &MacroState{state: s}
	s.Macros[name] = state
	return state
}

func NewState() (*State, error) {
	state := &State{configFile: configFile{path: "state.yaml", format: YAML}, Macros: make(map[string]*MacroState)}
	if err := state.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro state")
	}
	return state, nil
}
