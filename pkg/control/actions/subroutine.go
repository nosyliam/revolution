package actions

import (
	. "github.com/nosyliam/revolution/pkg/control/common"
)

type subroutineAction struct {
	name string
}

func (a *subroutineAction) Execute(deps *Dependencies) error {
	return deps.Exec(a.name)
}

func Subroutine(name string) Action {
	return &subroutineAction{name: name}
}
