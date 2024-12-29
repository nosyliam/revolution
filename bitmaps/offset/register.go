package offset

import (
	_ "embed"
	"github.com/nosyliam/revolution/bitmaps/registry"
)

func Register(registry registry.Registry) {
	registry.RegisterPng("tophoney", tophoney)
}

//go:embed tophoney.png
var tophoney []byte
