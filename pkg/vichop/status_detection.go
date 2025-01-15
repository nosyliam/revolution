package vichop

import (
	"fmt"
	"github.com/nosyliam/revolution/bitmaps"
	"github.com/nosyliam/revolution/pkg/common"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"image/color"
	"math"
)

type StatusDetector struct {
	manager *Manager
	macro   *common.Macro
	active  bool
}

func NewStatusDetector(manager *Manager, macro *common.Macro) *StatusDetector {
	return &StatusDetector{
		manager: manager,
		macro:   macro,
	}
}

func (s *StatusDetector) Tick() {
	screenshot := s.macro.GetWindow().Screenshot()
	if screenshot == nil {
		return
	}
	width, height := screenshot.Bounds().Dx(), screenshot.Bounds().Dy()

	// To detect Vicious bee attacking events, we'll use a heuristic method which involves
	// detecting the message's black background color up to an exact height, as well as the color of the
	// red text for at least 100 pixels. For defeated events, we can use template matching
	if s.active {
		for _, variant := range []string{"vicdefeated", "giftedvicdefeated"} {
			result, err := revimg.ImageSearch(bitmaps.Registry.Get(variant), screenshot, &revimg.SearchOptions{
				BoundStart: &revimg.Point{X: width - 355, Y: height - 200},
				BoundEnd:   &revimg.Point{X: width - 16},
				Variation:  1,
			})
			if err != nil {
				panic(err)
			}
			if len(result) > 0 {
				fmt.Println("Vicious bee defeated!")
				s.active = false
			}
		}
		return
	}

	// Start from the middle and expand search outwards (we'll only search for 150 pixels)
	var boxDetected = false
	for i := 1; i < 150; i++ {
		var dir = -1
		if i%2 == 0 {
			dir = 1
		}
		x := width - 186 - (int(math.Floor(float64(i/2))) * dir)
		var boxColor *color.RGBA
		var boxScan = 0
		var boxY = 0
		for y := height; y > height-100; y-- {
			if !boxDetected {
				rgba := screenshot.RGBAAt(x, y)
				if boxColor != nil && (rgba.R == boxColor.R || rgba.G == boxColor.G || rgba.B == boxColor.B) {
					boxScan++
				} else if rgba.R <= 14 && rgba.G <= 14 && rgba.B <= 14 && rgba.A == 255 {
					boxScan = 1
					boxY = y
					boxColor = &rgba
				} else {
					boxScan = 0
					boxColor = nil
				}
				if boxScan == 20 {
					boxDetected = true
					break
				}
			}
		}
		if boxDetected {
			var redPixels = 0
			for bI := 1; i < 250; i++ {
				var bDir = -1
				if i%2 == 0 {
					bDir = 1
				}
				bX := width - 186 - (int(math.Floor(float64(bI/2))) * bDir)
				for y := boxY; y > boxY-20; y-- {
					rgba := screenshot.RGBAAt(bX, y)
					if rgba.R > 200 {
						redPixels++
					}
				}
			}
			if redPixels > 100 {
				fmt.Println("Vicious bee attack detected")
				s.active = true
				return
			}
		}
	}
}
