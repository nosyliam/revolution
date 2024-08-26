package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type VariableName string

const (
	RetryCount      VariableName = "retry-count"
	UsePublicServer VariableName = "use-public-server"
)

func S(path string) interface{} {
	return nil
}

func V(name VariableName) func(macro *common.Macro) interface{} {
	return func(macro *common.Macro) interface{} {
		return macro.Scratch.Get(string(name))
	}
}

func Index(depth ...int) func(macro *common.Macro) interface{} {
	if len(depth) > 1 {
		panic("too many arguments")
	}
	var depthV int
	if len(depth) == 1 {
		depthV = depth[0]
	}
	return func(macro *common.Macro) interface{} {
		return macro.Scratch.LoopState.Index[depthV]
	}
}

func Window(macro *common.Macro) interface{} {
	return macro.Window
}

func LastError(macro *common.Macro) interface{} {
	return macro.Scratch.LastError
}
