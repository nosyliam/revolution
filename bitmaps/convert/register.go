package convert

import (
	_ "embed"
	"github.com/nosyliam/revolution/bitmaps/registry"
)

func Register(registry registry.Registry) {
	registry.RegisterPng("collectpollen", collectpollen)
}

//go:embed collectpollen.png
var collectpollen []byte
