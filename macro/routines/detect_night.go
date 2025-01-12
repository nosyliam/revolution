package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
)

const DetectNightRoutineKind RoutineKind = "detect-night"

var DetectNightRoutine = Actions{}

func init() {
	DetectNightRoutine.Register(DetectNightRoutineKind)
}
