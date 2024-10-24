package image

import "github.com/nosyliam/revolution/pkg/common"

type Modifier func(params ...interface{}) Value
type Value func(ctx *imageSearchContext) interface{}

type Step func(ctx *imageSearchContext)

type imageSearchContext struct {
	search *imageSearch
	macro  *common.Macro

	result     map[string][]Point
	lastResult []Point

	subX, subY     int
	x1, y1, x2, y2 int
	variance       int
	bitmap         string
}

func (i *imageSearchContext) execute() {

}

func NewImageSearchContext(search *imageSearch, macro *common.Macro) *imageSearchContext {
	return &imageSearchContext{
		search: search,
		macro:  macro,
	}
}

type imageSearch struct {
	steps []Step
}

func (i *imageSearch) Found() func(macro *common.Macro) bool {
	return func(macro *common.Macro) bool {
		ctx := NewImageSearchContext(i, macro)
		ctx.execute()
		return len(ctx.lastResult) > 0
	}
}

func (i *imageSearch) NotFound() func(macro *common.Macro) bool {
	return func(macro *common.Macro) bool {
		ctx := NewImageSearchContext(i, macro)
		ctx.execute()
		return len(ctx.lastResult) == 0
	}
}

func (i *imageSearch) Instances() func(macro *common.Macro) int {
	return func(macro *common.Macro) int {
		ctx := NewImageSearchContext(i, macro)
		ctx.execute()
		return len(ctx.lastResult)
	}
}

func Image(steps ...Step) *imageSearch {
	return &imageSearch{steps: steps}
}
