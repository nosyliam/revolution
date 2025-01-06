package movement

import (
	"context"
	"fmt"
	"github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

var mappedDefaultFunctions = make(map[string]bool)

func NoopFunction(L *lua.LState) int {
	return 0
}

func registerDirection(L *lua.LState) {
	newDirectionUD := func(d Direction) *lua.LUserData {
		ud := L.NewUserData()
		ud.Value = d
		L.SetMetatable(ud, L.GetTypeMetatable("DirectionMT"))
		return ud
	}

	mt := L.NewTypeMetatable("DirectionMT")
	L.SetField(mt, "__metatable", lua.LString("protected"))

	dirTable := L.NewTable()
	L.SetField(dirTable, "Forward", newDirectionUD(Forward))
	L.SetField(dirTable, "Backward", newDirectionUD(Backward))
	L.SetField(dirTable, "Left", newDirectionUD(Left))
	L.SetField(dirTable, "Right", newDirectionUD(Right))

	L.SetMetatable(dirTable, mt)
	L.SetGlobal("Direction", dirTable)
}

func registerKeys(L *lua.LState) {
	newKeyUD := func(k common.Key) *lua.LUserData {
		ud := L.NewUserData()
		ud.Value = k
		L.SetMetatable(ud, L.GetTypeMetatable("KeyMT"))
		return ud
	}

	mt := L.NewTypeMetatable("KeyMT")
	L.SetField(mt, "__metatable", lua.LString("protected"))

	dirTable := L.NewTable()
	L.SetField(dirTable, "Forward", newKeyUD(common.Forward))
	L.SetField(dirTable, "Backward", newKeyUD(common.Backward))
	L.SetField(dirTable, "Left", newKeyUD(common.Left))
	L.SetField(dirTable, "Right", newKeyUD(common.Right))
	L.SetField(dirTable, "RotLeft", newKeyUD(common.RotLeft))
	L.SetField(dirTable, "RotRight", newKeyUD(common.RotRight))
	L.SetField(dirTable, "ZoomIn", newKeyUD(common.ZoomIn))
	L.SetField(dirTable, "ZoomOut", newKeyUD(common.ZoomOut))
	L.SetField(dirTable, "E", newKeyUD(common.E))
	L.SetField(dirTable, "R", newKeyUD(common.R))
	L.SetField(dirTable, "L", newKeyUD(common.L))
	L.SetField(dirTable, "Esc", newKeyUD(common.Esc))
	L.SetField(dirTable, "Enter", newKeyUD(common.Enter))
	L.SetField(dirTable, "LShift", newKeyUD(common.LShift))
	L.SetField(dirTable, "Space", newKeyUD(common.Space))
	L.SetField(dirTable, "One", newKeyUD(common.One))
	L.SetField(dirTable, "Two", newKeyUD(common.Two))
	L.SetField(dirTable, "Three", newKeyUD(common.Three))
	L.SetField(dirTable, "Four", newKeyUD(common.Four))
	L.SetField(dirTable, "Five", newKeyUD(common.Five))
	L.SetField(dirTable, "Six", newKeyUD(common.Six))
	L.SetField(dirTable, "Seven", newKeyUD(common.Seven))

	L.SetMetatable(dirTable, mt)
	L.SetGlobal("Key", dirTable)
}

func registerPattern(L *lua.LState, meta *config.PatternMetadata) {
	metaTable := L.NewTable()
	L.SetField(metaTable, "Length", lua.LNumber(meta.Length))
	L.SetField(metaTable, "Width", lua.LNumber(meta.Width))
	L.SetField(metaTable, "Distance", lua.LNumber(meta.Width))

	mt := L.NewTypeMetatable("MetaMT")
	L.SetField(mt, "__metatable", lua.LString("protected"))
	L.SetMetatable(metaTable, mt)

	L.SetGlobal("Pattern", metaTable)
}

func LuaSleep(L *lua.LState) int {
	ms := L.CheckNumber(1)
	macro := L.Context().Value("macro").(*common.Macro)
	Sleep(int(ms), macro)
	return 0
}

func LuaWalk(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	dirUD := L.CheckUserData(1)
	studs := L.CheckNumber(2)
	if dir, ok := dirUD.Value.(Direction); !ok {
		L.RaiseError("expected a direction as the first argument, got %s", L.Get(1).Type().String())
		return 0
	} else {
		Walk(dir, float64(studs), macro)
	}
	return 0
}

func LuaWalkAsync(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	dirUD := L.CheckUserData(1)
	studs := L.CheckNumber(2)
	if dir, ok := dirUD.Value.(Direction); !ok {
		L.RaiseError("expected a direction as the first argument, got %s", L.Get(1).Type().String())
		return 0
	} else {
		WalkAsync(dir, float64(studs), macro)
	}
	return 0
}

func LuaKeyUp(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	keyUD := L.CheckUserData(1)
	if key, ok := keyUD.Value.(common.Key); !ok {
		L.RaiseError("expected a key as the first argument, got %s", L.Get(1).Type().String())
		return 0
	} else {
		<-macro.EventBus.KeyUp(macro, key)
	}
	return 0
}

func LuaKeyDown(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	keyUD := L.CheckUserData(1)
	if key, ok := keyUD.Value.(common.Key); !ok {
		L.RaiseError("expected a key as the first argument, got %s", L.Get(1).Type().String())
		return 0
	} else {
		<-macro.EventBus.KeyDown(macro, key)
	}
	return 0
}

func LuaKeyPress(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	keyUD := L.CheckUserData(1)
	if key, ok := keyUD.Value.(common.Key); !ok {
		L.RaiseError("expected a key as the first argument, got %s", L.Get(1).Type().String())
		return 0
	} else {
		<-macro.EventBus.KeyDown(macro, key)
		Sleep(50, macro) // TODO: Custom key delay
		<-macro.EventBus.KeyUp(macro, key)
	}
	return 0
}

func LuaQueryState(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	path := L.CheckString(1)
	if value, err := macro.MacroState.GetPath(path); err != nil {
		L.RaiseError("failed to fetch path: %v", err)
		return 0
	} else {
		switch v := value.(type) {
		case int:
			L.Push(lua.LNumber(v))
		case bool:
			L.Push(lua.LBool(v))
		case string:
			L.Push(lua.LString(v))
		case float64:
			L.Push(lua.LNumber(v))
		default:
			L.RaiseError(fmt.Sprintf("unexpected data type received from query: %T", v))
		}
	}
	return 1
}

func LuaQuerySetting(L *lua.LState) int {
	macro := L.Context().Value("macro").(*common.Macro)
	path := L.CheckString(1)
	if value, err := macro.Settings.GetPath(path); err != nil {
		L.RaiseError("failed to fetch path: %v", err)
		return 0
	} else {
		switch v := value.(type) {
		case int:
			L.Push(lua.LNumber(v))
		case bool:
			L.Push(lua.LBool(v))
		case string:
			L.Push(lua.LString(v))
		case float64:
			L.Push(lua.LNumber(v))
		default:
			L.RaiseError(fmt.Sprintf("unexpected data type received from query: %T", v))
		}
	}
	return 1
}

func ExecutePattern(pattern *Pattern, meta *config.PatternMetadata, macro *common.Macro) {
	L := lua.NewState()
	defer L.Close()
	registerDirection(L)
	registerKeys(L)
	registerPattern(L, meta)
	for name := range mappedDefaultFunctions {
		L.SetGlobal(fmt.Sprintf("Set%s", name), L.NewFunction(NoopFunction))
	}
	L.SetGlobal("SetName", L.NewFunction(NoopFunction))
	L.SetGlobal("Sleep", L.NewFunction(LuaSleep))
	L.SetGlobal("Walk", L.NewFunction(LuaWalk))
	L.SetGlobal("WalkAsync", L.NewFunction(LuaWalkAsync))
	L.SetGlobal("KeyUp", L.NewFunction(LuaKeyUp))
	L.SetGlobal("KeyDown", L.NewFunction(LuaKeyDown))
	L.SetGlobal("KeyPress", L.NewFunction(LuaKeyPress))
	L.SetGlobal("QueryState", L.NewFunction(LuaQueryState))
	L.SetGlobal("QuerySetting", L.NewFunction(LuaQuerySetting))
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "macro", macro)
	L.SetContext(ctx)
	go func() {
		watch := macro.Watch()
		defer macro.Unwatch(watch)
		select {
		case <-ctx.Done():
			return
		case resume := <-watch:
			if resume == nil {
				cancel()
			}
		}
	}()
	lf := L.NewFunctionFromProto(pattern.Proto)
	L.Push(lf)
	if err := L.PCall(0, 0, nil); err != nil {
		fmt.Println(err)
		macro.Status("Pattern execution error!")
		macro.Logger.LogDiscord(logging.Error, fmt.Sprintf("Failed to execute pattern %s: %v", pattern.Path, err), nil, nil)
		Sleep(5000, macro)
	}
}

func init() {
	t := reflect.TypeOf(config.PatternMetadata{})
	for i := 0; i < t.NumField(); i++ {
		meta := t.Field(i)
		mappedDefaultFunctions[meta.Name] = true
	}
}
