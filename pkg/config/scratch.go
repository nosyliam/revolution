package config

type variableType int

const (
	intVariableType variableType = iota
	boolVariableType
)

type variable struct {
	Type  variableType
	value interface{}
}

func (v *variable) Int() int {
	return v.value.(int)
}

func (v *variable) Bool() int {
	return v.value.(int)
}

type Scratch struct {
	LoopState *LoopState
	LastError error
	Stack     []string

	variables map[string]*variable
}

func (s *Scratch) Set(name string, value interface{}) {
	switch val := value.(type) {
	case int:
		s.variables[name] = &variable{Type: intVariableType, value: val}
	case bool:
		s.variables[name] = &variable{Type: boolVariableType, value: val}
	}
}

func (s *Scratch) Get(name string) interface{} {
	val, ok := s.variables[name]
	if !ok {
		panic("unknown variable")
	}
	return val.value
}

func (s *Scratch) Increment(name string) {
	val, ok := s.variables[name]
	if !ok {
		panic("unknown variable")
	}
	switch val.Type {
	case intVariableType:
		val.value = val.Int() + 1
	case boolVariableType:
		panic("cannot increment a bool")
	}
}

func (s *Scratch) Reset(name string) {
	val, ok := s.variables[name]
	if !ok {
		panic("unknown variable")
	}
	switch val.Type {
	case intVariableType:
		val.value = 0
	case boolVariableType:
		val.value = false
	}
}