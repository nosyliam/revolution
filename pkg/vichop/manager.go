package vichop

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
)

type Manager struct {
	Dataset *Dataset
	macros  []*common.Macro
}

func NewManager(state *config.Object[config.State]) *Manager {
	return &Manager{
		Dataset: NewDataset(state),
	}
}

func RegisterMacro(macro *common.Macro) {

}

func UnregisterMacro(macro *common.Macro) {

}
