package common

import revimg "github.com/nosyliam/revolution/pkg/image"

type LogicResult int

const (
	ContinueLogicResult LogicResult = iota
	RetryLogicResult
	ErrorLogicResult
)

type ActionResults struct {
	ImageSearchPoints []revimg.Point
	RetryCount        int
}

type Action interface {
	Execute(macro *Macro) error
}
