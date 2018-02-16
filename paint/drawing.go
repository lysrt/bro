package paint

import (
	"image"
	"image/color"

	"github.com/lysrt/bro/css"
)

type Canvas struct {
	*image.NRGBA
	color color.NRGBA
}

func NewCanvas(img *image.NRGBA) *Canvas {
	return &Canvas{NRGBA: img}
}

func (c *Canvas) SetColor(col css.Color) {
	c.color = color.NRGBA{
		A: col.A,
		R: col.R,
		G: col.G,
		B: col.B,
	}
}

// HLine draws a horizontal line
func (c *Canvas) HLine(x1, y, x2 int) {
	for ; x1 <= x2; x1++ {
		c.Set(x1, y, c.color)
	}
}

// VLine draws a veritcal line
func (c *Canvas) VLine(x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		c.Set(x, y1, c.color)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func (c *Canvas) Rect(x1, y1, x2, y2 int) {
	for ; y1 <= y2; y1++ {
		c.HLine(x1, y1, x2)
	}
}
