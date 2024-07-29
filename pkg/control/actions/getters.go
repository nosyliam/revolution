package actions

import "github.com/nosyliam/revolution/pkg/common"

func RetryCount(macro *common.Macro) int {
	return macro.Results.RetryCount
}

func Index(macro *common.Macro) int {
	return macro.Results.Index
}

func Window(macro *common.Macro) interface{} {
	return macro.Window
}

func True(*common.Macro) bool {
	return true
}

func False(*common.Macro) bool {
	return false
}
