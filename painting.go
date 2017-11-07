package main

import "image"

func Paint(layout *LayoutBox) (image.Image, error) {
	width, height := 50, 50
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	return img, nil
}
