package main

import (
	"fmt"
	"github.com/nosyliam/revolution/platform"
)

func main() {
	frames, _ := platform.WindowBackend.DisplayFrames()
	fmt.Println(frames)
}
