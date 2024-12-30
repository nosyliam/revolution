package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
)

type VariableName string

const (
	RetryCount      VariableName = "retry-count"
	UsePublicServer VariableName = "use-public-server"
	RestartSleep    VariableName = "restart-sleep"
	Offset          VariableName = "offset"
	OffsetX         VariableName = "offset-x"
	OffsetY         VariableName = "offset-y"
)

func S(path string) interface{} {
	return nil
}

// Get a variable as a concrete type
func V[T any](name VariableName) func(macro *common.Macro) T {
	return func(macro *common.Macro) T {
		return macro.Scratch.Get(string(name)).(T)
	}
}

// Get a variable as an interface
func VI(name VariableName) func(macro *common.Macro) interface{} {
	return func(macro *common.Macro) interface{} {
		return macro.Scratch.Get(string(name))
	}
}

// Get a setting path as a concrete type
func P[T any](path string) func(macro *common.Macro) T {
	return func(macro *common.Macro) T {
		return *config.Concrete[T](macro.Settings, path)
	}
}

func Index(depth ...int) func(macro *common.Macro) int {
	if len(depth) > 1 {
		panic("too many arguments")
	}
	var depthV = 0
	if len(depth) == 1 {
		depthV = depth[0]
	}
	return func(macro *common.Macro) int {
		return macro.Scratch.LoopState.Index[depthV]
	}
}

func Window(macro *common.Macro) interface{} {
	return macro.Root.Window
}

func LastError(macro *common.Macro) interface{} {
	return macro.Scratch.LastError
}
