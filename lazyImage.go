package lazyImage

import (
	"image"
)

type LazyImage struct {
	TheColorModel image.ColorModel
	TheBounds     image.Rectangle
	Fill          image.Color
	lastPos       int
}

func (t *LazyImage) ColorModel() image.ColorModel {
	return t.TheColorModel
}

func (t *LazyImage) Bounds() image.Rectangle {
	return t.TheBounds
}

func (t *LazyImage) At(x, y int) image.Color {
	t.lastPos += 1
	//pos:=x+y*t.TheBounds.Dx()
	if x == y && x%100 == 0 {
		print(t.lastPos)
		print(" _ ")
		print(x)
		print(", ")
		print(y, "\n")
		//t.lastPos=pos
	}
	c := new(image.NRGBAColor)
	c.R = uint8((y/1 + x/2) % 256)
	c.A = 255
	return c
}


func NewLazyImage(width, height int) *LazyImage {
	fill := image.Black
	bounds := image.Rect(0, 0, width, height)
	return &LazyImage{image.NRGBAColorModel, bounds, fill, -1}
}
