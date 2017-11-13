package main

import (
	"image"
	"image/draw"

	"github.com/lysrt/bro/css"
)

type DisplayList []DisplayCommand

type DisplayCommand interface {
	paint(*Canvas)
}

type SolidColor struct {
	color css.Color
	rect  Rect
}

func (c *SolidColor) paint(img *Canvas) {
	img.SetColor(c.color)
	x0 := int(c.rect.x)
	y0 := int(c.rect.y)
	x1 := int(c.rect.x + c.rect.width)
	y1 := int(c.rect.y + c.rect.height)
	img.Rect(x0, y0, x1, y1)
}

func Paint(layoutRoot *LayoutBox) (image.Image, error) {
	displayList := buildDisplayList(layoutRoot)

	width := int(layoutRoot.dimensions.content.width)
	height := int(layoutRoot.dimensions.content.height)
	// width, height := 500, 500
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	canvas := NewCanvas(img)
	for _, item := range displayList {
		paintItem(canvas, item)
	}

	return img, nil
}

func buildDisplayList(layoutRoot *LayoutBox) DisplayList {
	var list DisplayList
	renderLayoutBox(&list, layoutRoot)
	return list
}

func renderLayoutBox(list *DisplayList, layoutBox *LayoutBox) {
	renderBackground(list, layoutBox)
	renderBorders(list, layoutBox)
	// TODO render text

	for _, child := range layoutBox.children {
		renderLayoutBox(list, child)
	}
}

func renderBackground(list *DisplayList, layoutBox *LayoutBox) {
	if color, ok := getColor(layoutBox, "background-color"); ok {
		*list = append(*list, &SolidColor{color: color, rect: layoutBox.dimensions.borderBox()})
	}
}

func renderBorders(list *DisplayList, layoutBox *LayoutBox) {
	var color css.Color
	var ok bool

	if color, ok = getColor(layoutBox, "border-color"); !ok {
		// TODO is it OK to return if no border-color is specified?
		return
	}

	d := layoutBox.dimensions
	borderBox := d.borderBox()

	// Left border
	*list = append(*list, &SolidColor{color: color,
		rect: Rect{
			x:      borderBox.x,
			y:      borderBox.y,
			width:  d.border.left,
			height: borderBox.height,
		},
	})

	// Right border
	*list = append(*list, &SolidColor{color: color,
		rect: Rect{
			x:      borderBox.x + borderBox.width - d.border.right,
			y:      borderBox.y,
			width:  d.border.right,
			height: borderBox.height,
		},
	})

	// Top border
	*list = append(*list, &SolidColor{color: color,
		rect: Rect{
			x:      borderBox.x,
			y:      borderBox.y,
			width:  borderBox.width,
			height: d.border.top,
		},
	})

	// Bottom border
	*list = append(*list, &SolidColor{color: color,
		rect: Rect{
			x:      borderBox.x,
			y:      borderBox.y + borderBox.height - d.border.bottom,
			width:  borderBox.width,
			height: d.border.bottom,
		},
	})
}

func getColor(layoutBox *LayoutBox, name string) (color css.Color, ok bool) {
	switch layoutBox.boxType {
	case BlockNode:
		fallthrough
	case InlineNode:
		value, ok := layoutBox.styledNode.Value(name)
		if !ok {
			return css.Color{}, false
		}
		return value.Color, true
	case AnonymousBlock:
		return css.Color{}, false
	default:
		panic("Unknown boxType")
	}
}

func paintItem(c *Canvas, command DisplayCommand) {
	switch t := command.(type) {
	case *SolidColor:
		t.paint(c)
	default:
		panic("Unexpect DisplayCommand type")
	}
}
