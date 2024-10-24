package image

import "github.com/nosyliam/revolution/pkg/common"

func X1(ctx *imageSearchContext) interface{} {
	return ctx.x1
}

func Y1(ctx *imageSearchContext) interface{} {
	return ctx.y1
}

func X2(ctx *imageSearchContext) interface{} {
	return ctx.x2
}

func Y2(ctx *imageSearchContext) interface{} {
	return ctx.y2
}

func XY1(ctx *imageSearchContext) interface{} {
	return Point{X: ctx.x1, Y: ctx.y1}
}

func XY2(ctx *imageSearchContext) interface{} {
	return Point{X: ctx.x2, Y: ctx.y2}
}

func Height(ctx *imageSearchContext) interface{} {
	return Point{X: ctx.x2, Y: ctx.y2}
}

// Add adds together two or more values
func Add(params ...interface{}) Value {
	if len(params) != 2 {
		panic("two parameters required for modifier: add")
	}
	return func(ctx *imageSearchContext) interface{} {
		var value int
		for _, param := range params {
			switch v := param.(type) {
			case int:
				value += param.(int)
			case Value:
				value += v(ctx).(int)
			case func(macro *common.Macro) int:
				value += v(ctx.macro)
			default:
				panic("invalid parameter type: add")
			}
		}
		return value
	}
}

// Set returns the value of the last parameter
func Set(params ...interface{}) Value {
	if len(params) == 0 {
		panic("one or more parameters required for modifier: set")
	}
	return func(ctx *imageSearchContext) interface{} {
		switch v := params[len(params)-1].(type) {
		case int:
			return v
		case Value:
			return v(ctx)
		case func(macro *common.Macro) int:
			return v(ctx.macro)
		default:
			panic("invalid parameter type: add")
		}
	}
}

// Sub subtracts the first parameter from one or more parameters
func Sub(params ...interface{}) Value {
	if len(params) != 2 {
		panic("two or more parameters required for modifier: sub")
	}
	return func(ctx *imageSearchContext) interface{} {
		var value int
		for i, param := range params {
			switch v := param.(type) {
			case int:
				if i == 0 {
					value = v
					continue
				}
				value -= v
			case Value:
				if i == 0 {
					value = v(ctx).(int)
					continue
				}
				value -= v(ctx).(int)
			case func(macro *common.Macro) int:
				if i == 0 {
					value = v(ctx.macro)
					continue
				}
				value -= v(ctx.macro)
			default:
				panic("invalid parameter type: sub")
			}
		}
		return value
	}
}

// Div divides (and floors) two parameters

// Mul multiplies two or more parameters

func Result(params ...interface{}) Value {
	if len(params) != 2 {
		panic("more than one parameter required for modifier: result")
	}
	return func(ctx *imageSearchContext) interface{} {
		var point Point
		if len(params) == 0 {
			point = ctx.lastResult[0]
		} else {
			point = ctx.result[params[0].(string)][0]
		}
		return point
	}
}

func SelectCoordinate(params ...interface{}) Step {
	return func(ctx *imageSearchContext) {
		switch fn := params[0].(type) {
		case Modifier:
			params = params[1:]
			if len(params) != 4 {
				panic("expected 4 parameters for coordinate")
			}
			var getters = []Value{X1, Y1, X2, Y2}
			var values []int
			for i, paramValue := range params {
				switch param := paramValue.(type) {
				case int:
					values = append(values, fn(getters[i], param)(ctx).(int))
				case func(macro *common.Macro) int:
					values = append(values, fn(getters[i], param(ctx.macro))(ctx).(int))
				case Value:
					values = append(values, fn(getters[i], param(ctx).(int))(ctx).(int))
				default:
					panic("invalid value type")
				}
			}
			ctx.x1, ctx.y1, ctx.x2, ctx.y2 = values[0], values[1], values[2], values[3]
		case Value:
			point := fn(ctx).(Point)
			ctx.x1, ctx.y1, ctx.x2, ctx.y2 = point.X, point.Y, point.X, point.Y
		}
	}
}

func singleCoordinateSelect(value Value, mod func(ctx *imageSearchContext, val int), params ...interface{}) Step {
	return func(ctx *imageSearchContext) {
		switch fn := params[0].(type) {
		case Modifier:
			params = params[1:]
			if len(params) != 4 {
				panic("expected one parameter for single coordinate selection")
			}
			switch param := params[0].(type) {
			case int:
				mod(ctx, fn(value, param)(ctx).(int))
			case Value:
				mod(ctx, fn(value, param(ctx).(int))(ctx).(int))
			case func(macro *common.Macro) int:
				mod(ctx, fn(value, param(ctx.macro))(ctx).(int))
			default:
				panic("invalid value type")
			}
		case Value:
			mod(ctx, fn(ctx).(int))
		case func(macro *common.Macro) int:
			mod(ctx, fn(ctx.macro))
		case int:
			mod(ctx, fn)
		}
	}
}

func SelectX1(params ...interface{}) Step {
	return singleCoordinateSelect(X1, func(ctx *imageSearchContext, val int) { ctx.x1 = val }, params...)
}

func SelectY1(params ...interface{}) Step {
	return singleCoordinateSelect(Y1, func(ctx *imageSearchContext, val int) { ctx.y1 = val }, params...)
}

func SelectX2(params ...interface{}) Step {
	return singleCoordinateSelect(X2, func(ctx *imageSearchContext, val int) { ctx.x2 = val }, params...)
}

func SelectY2(params ...interface{}) Step {
	return singleCoordinateSelect(Y2, func(ctx *imageSearchContext, val int) { ctx.y2 = val }, params...)
}
