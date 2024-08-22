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

type MacroState struct {
	AccountName string `yaml:"accountName"`
	Status      string `state:"status"`

	LoopState *LoopState
	LastError error
	Stack     []string

	PrivateServerAttempts int
	UsePublicServer       bool

	ClaimedHive *int
}

func (m *MacroState) Key() string {
	return m.AccountName
}

type State struct {
	ActiveAccount *string          `yaml:"activeAccount"`
	Macros        List[MacroState] `yaml:"macros"`
}

func (s *State) State(name string) *MacroState {
	/*for _, macro := range s.Macros.data {
		if macro.AccountName == name {
			return macro
		}
	}
	if err := s.AppendPath(fmt.Sprintf("macros[%s]", name)); err != nil {
		panic(err)
	}*/
	return s.State(name)
}

func NewState() (*Object[State], error) {
	state := File[State]{path: "state.yaml", format: YAML}
	if err := state.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load macro state")
	}
	return state.Object(), nil
}
