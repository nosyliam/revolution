package main

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/nosyliam/revolution/platform"
)

func main() {
	fmt.Println(platform.WindowBackend.DisplayFrames())
	id, _ := platform.WindowBackend.OpenWindow(window.JoinOptions{})
	platform.WindowBackend.Screenshot(id)
}
