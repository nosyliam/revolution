package actions

import (
	. "github.com/nosyliam/revolution/pkg/common"
)

type subroutineAction struct {
	name string
}

func (a *subroutineAction) Execute(deps *Macro) error {
	return deps.ExecRoutine(a.name)
}

func Subroutine(name string) Action {
	return &subroutineAction{name: name}
}
