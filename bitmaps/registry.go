package bitmaps

import (
	"bytes"
	"github.com/nosyliam/revolution/bitmaps/convert"
	"github.com/nosyliam/revolution/bitmaps/offset"
	"image"
	"image/draw"
	"image/png"

	_ "github.com/nosyliam/revolution/bitmaps/offset"
)

var Registry = &bitmapRegistry{bitmaps: make(map[string]*image.RGBA)}

type bitmapRegistry struct {
	bitmaps map[string]*image.RGBA
}

func (b *bitmapRegistry) Get(name string) *image.RGBA {
	return b.bitmaps[name]
}

func (b *bitmapRegistry) RegisterPng(name string, data []byte) {
	reader := bytes.NewReader(data)
	img, err := png.Decode(reader)
	if err != nil {
		panic(err)
	}
	rgba, ok := img.(*image.RGBA)
	if !ok {
		b := img.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	}
	b.bitmaps[name] = rgba
}

func (b *bitmapRegistry) initialize() {
	offset.Register(b)
	convert.Register(b)
}

func init() {
	Registry.initialize()
}
