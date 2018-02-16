package paint

import (
	"image"
	"image/draw"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/layout"
)

type DisplayList []DisplayCommand

type DisplayCommand interface {
	paint(*Canvas)
}

type SolidColor struct {
	color css.Color
	rect  layout.Rect
}

func (c *SolidColor) paint(img *Canvas) {
	img.SetColor(c.color)
	x0 := int(c.rect.X)
	y0 := int(c.rect.Y)
	x1 := int(c.rect.X + c.rect.Width)
	y1 := int(c.rect.Y + c.rect.Height)
	img.Rect(x0, y0, x1, y1)
}

func Paint(layoutRoot *layout.LayoutBox) (image.Image, error) {
	displayList := buildDisplayList(layoutRoot)

	width := int(layoutRoot.Dimensions.Content.Width)
	height := int(layoutRoot.Dimensions.Content.Height)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	canvas := NewCanvas(img)
	for _, item := range displayList {
		paintItem(canvas, item)
	}

	return img, nil
}

func buildDisplayList(layoutRoot *layout.LayoutBox) DisplayList {
	var list DisplayList
	renderLayoutBox(&list, layoutRoot)
	return list
}

func renderLayoutBox(list *DisplayList, layoutBox *layout.LayoutBox) {
	renderBackground(list, layoutBox)
	renderBorders(list, layoutBox)
	// TODO render text

	for _, child := range layoutBox.Children {
		renderLayoutBox(list, child)
	}
}

func renderBackground(list *DisplayList, layoutBox *layout.LayoutBox) {
	if color, ok := getColor(layoutBox, "background-color"); ok {
		*list = append(*list, &SolidColor{color: color, rect: layoutBox.Dimensions.BorderBox()})
	}
}

func renderBorders(list *DisplayList, layoutBox *layout.LayoutBox) {
	var color css.Color
	var ok bool

	if color, ok = getColor(layoutBox, "border-color"); !ok {
		// TODO is it OK to return if no border-color is specified?
		return
	}

	d := layoutBox.Dimensions
	borderBox := d.BorderBox()

	// Left border
	*list = append(*list, &SolidColor{color: color,
		rect: layout.Rect{
			X:      borderBox.X,
			Y:      borderBox.Y,
			Width:  d.Border.Left,
			Height: borderBox.Height,
		},
	})

	// Right border
	*list = append(*list, &SolidColor{color: color,
		rect: layout.Rect{
			X:      borderBox.X + borderBox.Width - d.Border.Right,
			Y:      borderBox.Y,
			Width:  d.Border.Right,
			Height: borderBox.Height,
		},
	})

	// Top border
	*list = append(*list, &SolidColor{color: color,
		rect: layout.Rect{
			X:      borderBox.X,
			Y:      borderBox.Y,
			Width:  borderBox.Width,
			Height: d.Border.Top,
		},
	})

	// Bottom border
	*list = append(*list, &SolidColor{color: color,
		rect: layout.Rect{
			X:      borderBox.X,
			Y:      borderBox.Y + borderBox.Height - d.Border.Bottom,
			Width:  borderBox.Width,
			Height: d.Border.Bottom,
		},
	})
}

func getColor(layoutBox *layout.LayoutBox, name string) (color css.Color, ok bool) {
	switch layoutBox.BoxType {
	case layout.BlockNode:
		fallthrough
	case layout.InlineNode:
		value, ok := layoutBox.StyledNode.Value(name)
		if !ok {
			return css.Color{}, false
		}
		return value.Color, true
	case layout.AnonymousBlock:
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
