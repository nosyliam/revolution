package actions

import "github.com/nosyliam/revolution/pkg/common"

type closeWindowAction struct{}

func (c *closeWindowAction) Execute(macro *common.Macro) error {
	if macro.Window != nil {
		return macro.Window.Close()
	}
}
