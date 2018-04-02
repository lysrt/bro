package paint

import (
	"image"

	"github.com/fogleman/gg"

	"github.com/lysrt/bro/css"
)

type Canvas struct {
	context *gg.Context
}

func NewCanvas(width, height int) *Canvas {
	context := gg.NewContext(width, height)
	return &Canvas{context: context}
}

func (c *Canvas) Image() image.Image {
	return c.context.Image()
}

func (c *Canvas) SetColor(col css.Color) {
	c.context.SetRGBA255(col.R, col.G, col.B, col.A)
}

// HLine draws a horizontal line
func (c *Canvas) HLine(x1, y, x2 int) {
	c.context.DrawLine(float64(x1), float64(y), float64(x2), float64(y))
	c.context.Stroke()
}

// VLine draws a veritcal line
func (c *Canvas) VLine(x, y1, y2 int) {
	c.context.DrawLine(float64(x), float64(y1), float64(x), float64(y2))
	c.context.Stroke()
}

// Rect draws a rectangle utilizing HLine() and VLine()
func (c *Canvas) Rect(x, y, width, height int) {
	c.context.DrawRectangle(float64(x), float64(y), float64(width), float64(height))
	c.context.Fill()
}
