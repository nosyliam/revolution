package actions

import "github.com/nosyliam/revolution/pkg/common"

type setVariableAction struct {
	name VariableName
	val  interface{}
}

func (a *setVariableAction) Execute(macro *common.Macro) error {
	macro.Scratch.Set(string(a.name), a.val)
	return nil
}

func Set(name VariableName, val interface{}) common.Action {
	return &setVariableAction{name: name, val: val}
}

type resetVariableAction struct {
	name VariableName
}

func (a *resetVariableAction) Execute(macro *common.Macro) error {
	macro.Scratch.Reset(string(a.name))
	return nil
}

func Reset(name VariableName) common.Action {
	return &resetVariableAction{name: name}
}

type incrementVariableAction struct {
	name VariableName
}

func (a *incrementVariableAction) Execute(macro *common.Macro) error {
	macro.Scratch.Increment(string(a.name))
	return nil
}

func Increment(name VariableName) common.Action {
	return &incrementVariableAction{name: name}
}
