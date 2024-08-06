package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

func RetryCount(macro *common.Macro) int {
	return macro.Results.RetryCount
}

func Index(depth ...int) func(macro *common.Macro) int {
	if len(depth) > 1 {
		panic("invalid arguments")
	}
	var depthV int
	if len(depth) == 1 {
		depthV = depth[0]
	}
	return func(macro *common.Macro) int {
		return macro.State.LoopState.Index[depthV]
	}
}

func Window(macro *common.Macro) interface{} {
	return macro.Window
}

func LastError(macro *common.Macro) interface{} {
	return macro.State.LastError
}
