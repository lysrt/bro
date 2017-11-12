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

type BoxType string

const (
	BlockNode      BoxType = "block"
	InlineNode             = "inline"
	AnonymousBlock         = "anon"
)

type LayoutBox struct {
	dimensions Dimensions
	boxType    BoxType
	styledNode *StyledNode
	children   []*LayoutBox
}

func newLayoutBox(boxType BoxType, styledNode *StyledNode) *LayoutBox {
	var children []*LayoutBox
	return &LayoutBox{
		boxType:    boxType,
		styledNode: styledNode,
		children:   children,
	}
}

func GenerateLayoutTree(styleTree *StyledNode) *LayoutBox {
	var boxType BoxType
	switch styleTree.Display() {
	case Inline:
		boxType = InlineNode
	case Block:
		boxType = BlockNode
	case None:
		panic("Root StyledNode has display:none")
	}

	root := newLayoutBox(boxType, styleTree)

	for _, child := range styleTree.Children {
		switch child.Display() {
		case Inline:
			ic := root.getInlineContainer()
			ic.children = append(ic.children, GenerateLayoutTree(child))
		case Block:
			root.children = append(root.children, GenerateLayoutTree(child))
		case None:
			// Skip
		}
	}

	return root
}

// If a block node contains an inline child:
func (box *LayoutBox) getInlineContainer() *LayoutBox {
	switch box.boxType {
	case InlineNode:
		fallthrough
	case AnonymousBlock:
		return box
	case BlockNode:
		// If we've just generated an anonymous block box, keep using it.
		// Otherwise, create a new one.
		lastChild := box.children[len(box.children)-1]
		switch lastChild.boxType {
		case AnonymousBlock:
			box.children = append(box.children, newLayoutBox(AnonymousBlock, nil))
		}
		return lastChild
	}
	panic("No more cases to switch")
}

func (box *LayoutBox) Layout(containingBlock Dimensions) {
	switch box.boxType {
	case InlineNode:
		// TODO
	case BlockNode:
		box.layoutBlock(containingBlock)
	case AnonymousBlock:
		// TODO
	}
}

func (box *LayoutBox) layoutBlock(containingBlock Dimensions) {
	// First go down the LayoutTree to compute the widths from parents' widths
	// Then go up the tree to compute heights form children's heights

	box.calculateWidth(containingBlock)

	box.calculatePosition(containingBlock)

	box.layoutBlockChildren()

	box.calculateHeight()
}
