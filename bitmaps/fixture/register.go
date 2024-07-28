package fixture

import (
	_ "embed"
	"github.com/nosyliam/revolution/bitmaps/registry"
)

func Register(registry registry.Registry) {
	registry.RegisterPng("honeyfixture", honeyFixture)
	registry.RegisterPng("honeyfixture2", honeyFixture2)
}

//go:embed honey.png
var honeyFixture []byte

//go:embed honey2.png
var honeyFixture2 []byte
