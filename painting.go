package main

import (
	"image"
	"image/draw"
)

func Paint(layout *LayoutBox) (image.Image, error) {
	width, height := 256, 256
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	return img, nil
}
