package main

import (
	"github.com/nosyliam/revolution/platform"
)

func main() {
	platform.ControlBackend.Sleep(10, make(chan struct{}))
}
