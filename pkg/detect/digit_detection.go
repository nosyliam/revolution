package detect

import (
	"github.com/nosyliam/revolution/bitmaps"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"github.com/pkg/errors"
	"image"
	"math"
)

var digits = make(map[*image.RGBA]int)

func generateNumber(digits []int) int {
	finalNumber := 0
	for i, digit := range digits {
		finalNumber += digit * int(math.Pow(10, float64(i)))
	}
	return finalNumber
}

// DetectDigits detects the digits inside an image. Returns -1 if no digits were found
func DetectDigits(image *image.RGBA) (int, error) {
	var detectedDigits []int
	if len(digits) == 0 {
		digits[bitmaps.Registry.Get("digit_zero")] = 0
		digits[bitmaps.Registry.Get("digit_one")] = 1
		digits[bitmaps.Registry.Get("digit_two")] = 2
		digits[bitmaps.Registry.Get("digit_three")] = 3
		digits[bitmaps.Registry.Get("digit_four")] = 4
		digits[bitmaps.Registry.Get("digit_five")] = 5
		digits[bitmaps.Registry.Get("digit_six")] = 6
		digits[bitmaps.Registry.Get("digit_seven")] = 7
		digits[bitmaps.Registry.Get("digit_eight")] = 8
		digits[bitmaps.Registry.Get("digit_nine")] = 9
	}
	searchX := 38
	for needle, digit := range digits {
		points, err := revimg.ImageSearch(needle, image,
			&revimg.SearchOptions{
				BoundStart:      &revimg.Point{Y: 20},
				BoundEnd:        &revimg.Point{X: searchX, Y: 38},
				SearchDirection: 8,
				Variation:       10,
			},
		)
		if err != nil {
			return -1, errors.Wrap(err, "failed to perform image search")
		}
		if len(points) == 0 {
			continue
		}
		searchX = points[0].X
		detectedDigits = append(detectedDigits, digit)
	}
	if len(detectedDigits) == 0 {
		return -1, nil
	}
	return generateNumber(detectedDigits), nil
}
