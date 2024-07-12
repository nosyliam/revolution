package offset

import (
	_ "embed"
	"github.com/nosyliam/revolution/bitmaps/registry"
)

func Register(registry registry.Registry) {
	registry.RegisterPng("toppollen", toppollen)
}

//go:embed toppollen.png
var toppollen []byte
