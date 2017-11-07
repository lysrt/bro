package main

type Dimensions struct {
	// Position of the content area relative to the document origin:
	content Rect

	// Surrounding edges:
	padding, border, margin EdgeSizes
}

type Rect struct {
	x, y, width, height float32
}

type EdgeSizes struct {
	left, right, top, bottom float32
}

type BoxType int

const (
	BlockNode BoxType = iota
	InlineNode
	AnonymousBlock
)

type LayoutBox struct {
	dimensions Dimensions
	boxType    BoxType
	children   []LayoutBox
}

func GenerateLayoutTree(style *StyledNode) (*LayoutBox, error) {
	return nil, nil
}
