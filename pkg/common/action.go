package common

import revimg "github.com/nosyliam/revolution/pkg/image"

type ActionResults struct {
	ImageSearchPoints []revimg.Point
	RetryCount        int
}

type Action interface {
	Execute(macro *Macro) error
}