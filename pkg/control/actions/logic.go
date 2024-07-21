package actions

import "github.com/nosyliam/revolution/pkg/common"

type logicAction struct {
	name string
}

func (a *logicAction) Execute(macro common.Macro) error {
	return macro.Exec(a.name)
}

func Logic(name string) common.Action {
	return &subroutineAction{name: name}
}
