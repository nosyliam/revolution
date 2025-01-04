package movement

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

var mappedDefaultFunctions = make(map[string]bool)

func NoopFunction(L *lua.LState) int {
	return 0
}

func GetLength(L *lua.LState) int {
	return 1
}

func GetWidth(L *lua.LState) int {
	return 1
}

func GetDistance(L *lua.LState) int {
	return 1
}

func ExecutePattern(pattern *Pattern, macro *common.Macro) error {
	L := lua.NewState()
	defer L.Close()
	for name, _ := range mappedDefaultFunctions {
		L.SetGlobal(fmt.Sprintf("Set%s", name), L.NewFunction(NoopFunction))
	}
	return nil
}

func init() {
	t := reflect.TypeOf(config.PatternMetadata{})
	for i := 0; i < t.NumField(); i++ {
		meta := t.Field(i)
		mappedDefaultFunctions[meta.Name] = true
	}
}
