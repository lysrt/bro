package main

import "github.com/lysrt/bro/css"

type Dimensions struct {
	// Position of the content area relative to the document origin:
	content Rect

	// Surrounding edges:
	padding, border, margin EdgeSizes
}

// marginBox returns the area covered by the content area plus padding, borders, and margin
func (d Dimensions) marginBox() Rect {
	return d.borderBox().expandedBy(d.margin)
}

// borderBox returns the area covered by the content area plus padding and borders
func (d Dimensions) borderBox() Rect {
	return d.paddingBox().expandedBy(d.border)
}

// paddingBox returns the area covered by the content area plus its padding
func (d Dimensions) paddingBox() Rect {
	return d.content.expandedBy(d.padding)
}

type Rect struct {
	x, y, width, height float32
}

func (r Rect) expandedBy(edge EdgeSizes) Rect {
	return Rect{
		x:      r.x - edge.left,
		y:      r.y - edge.top,
		width:  r.width + edge.left + edge.right,
		height: r.height + edge.top + edge.bottom,
	}
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

func (box *LayoutBox) calculateWidth(containingBlock Dimensions) {
	style := box.styledNode

	// width has initial value auto
	auto := css.Value{Keyword: "auto"}
	width, ok := style.Value("width")
	if !ok {
		width = auto
	}

	// margin, border, and padding have initial value 0
	zero := css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}

	var marginLeft, marginRight, paddingLeft, paddingRight, borderLeft, borderRight css.Value

	if marginLeft, ok = style.Value("margin-left"); !ok {
		if marginLeft, ok = style.Value("margin"); !ok {
			marginLeft = zero
		}
	}
	if marginRight, ok = style.Value("margin-right"); !ok {
		if marginRight, ok = style.Value("margin"); !ok {
			marginRight = zero
		}
	}

	if borderLeft, ok = style.Value("border-left-width"); !ok {
		if borderLeft, ok = style.Value("border-width"); !ok {
			borderLeft = zero
		}
	}
	if borderRight, ok = style.Value("border-right-width"); !ok {
		if borderRight, ok = style.Value("border-width"); !ok {
			borderRight = zero
		}
	}

	if paddingLeft, ok = style.Value("padding-left"); !ok {
		if paddingLeft, ok = style.Value("padding"); !ok {
			paddingLeft = zero
		}
	}
	if paddingRight, ok = style.Value("padding-right"); !ok {
		if paddingRight, ok = style.Value("padding"); !ok {
			paddingRight = zero
		}
	}

	// Formula for block width: https://www.w3.org/TR/CSS2/visudet.html#blockwidth
	// Auto must count as zero
	total := marginLeft.ToPx() + marginRight.ToPx() + borderLeft.ToPx() + borderRight.ToPx() +
		paddingLeft.ToPx() + paddingRight.ToPx() + width.ToPx()

	// Checking if the box is too big
	// If width is not auto and the total is wider than the container, treat auto margins as 0.
	if width != auto && total > containingBlock.content.width {
		if marginLeft == auto {
			marginLeft = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}
		if marginRight == auto {
			marginRight = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}
	}

	// Check for over or underflow, and adjust "auto" dimensions accordingly
	underflow := containingBlock.content.width - total

	widthAuto := width == auto
	marginLeftAuto := marginLeft == auto
	marginRightAuto := marginRight == auto

	if !widthAuto && !marginLeftAuto && !marginRightAuto {
		// If the values are overconstrained, calculate margin_right
		marginRight = css.Value{Length: css.Length{Quantity: marginRight.ToPx() + underflow, Unit: css.Px}}
	} else if !widthAuto && !marginLeftAuto && marginRightAuto {
		// If exactly one size is auto, its used value follows from the equality
		marginRight = css.Value{Length: css.Length{Quantity: underflow, Unit: css.Px}}
	} else if !widthAuto && marginLeftAuto && !marginRightAuto {
		// Idem
		marginLeft = css.Value{Length: css.Length{Quantity: underflow, Unit: css.Px}}
	} else if widthAuto {
		// If width is set to auto, any other auto values become 0
		if marginLeft == auto {
			marginLeft = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}
		if marginRight == auto {
			marginRight = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}

		if underflow >= 0.0 {
			// Expand width to fill the underflow
			width = css.Value{Length: css.Length{Quantity: underflow, Unit: css.Px}}
		} else {
			// Width can't be negative. Adjust the right margin instead
			width = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
			marginRight = css.Value{Length: css.Length{Quantity: marginRight.ToPx() + underflow, Unit: css.Px}}
		}

	} else if !widthAuto && marginLeftAuto && marginRightAuto {
		// If margin-left and margin-right are both auto, their used values are equal
		marginLeft = css.Value{Length: css.Length{Quantity: underflow / 2.0, Unit: css.Px}}
		marginRight = css.Value{Length: css.Length{Quantity: underflow / 2.0, Unit: css.Px}}
	}

	box.dimensions.content.width = width.ToPx()

	box.dimensions.padding.left = paddingLeft.ToPx()
	box.dimensions.padding.right = paddingRight.ToPx()

	box.dimensions.border.left = borderLeft.ToPx()
	box.dimensions.border.right = borderRight.ToPx()

	box.dimensions.margin.left = marginLeft.ToPx()
	box.dimensions.margin.right = marginRight.ToPx()
}

func (box *LayoutBox) calculatePosition(containingBlock Dimensions) {
	style := box.styledNode

	// margin, border, and padding have initial value 0
	zero := css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}

	var marginTop, marginBottom, borderTop, borderBottom, paddingTop, paddingBottom css.Value
	var ok bool

	if marginTop, ok = style.Value("margin-top"); !ok {
		if marginTop, ok = style.Value("margin"); !ok {
			marginTop = zero
		}
	}
	if marginBottom, ok = style.Value("margin-bottom"); !ok {
		if marginBottom, ok = style.Value("margin"); !ok {
			marginBottom = zero
		}
	}

	if borderTop, ok = style.Value("border-top-width"); !ok {
		if borderTop, ok = style.Value("border-width"); !ok {
			borderTop = zero
		}
	}
	if borderBottom, ok = style.Value("border-bottom-width"); !ok {
		if borderBottom, ok = style.Value("border-width"); !ok {
			borderBottom = zero
		}
	}

	if paddingTop, ok = style.Value("padding-top"); !ok {
		if paddingTop, ok = style.Value("padding"); !ok {
			paddingTop = zero
		}
	}
	if paddingBottom, ok = style.Value("padding-bottom"); !ok {
		if paddingBottom, ok = style.Value("padding"); !ok {
			paddingBottom = zero
		}
	}

	box.dimensions.margin.top = marginTop.ToPx()
	box.dimensions.margin.bottom = marginBottom.ToPx()
	box.dimensions.border.top = borderTop.ToPx()
	box.dimensions.border.bottom = borderBottom.ToPx()
	box.dimensions.padding.top = paddingTop.ToPx()
	box.dimensions.padding.bottom = paddingBottom.ToPx()

	box.dimensions.content.x = containingBlock.content.x +
		box.dimensions.margin.left + box.dimensions.border.left + box.dimensions.padding.left

	// Position the box below all the previous boxes in the container.
	// Making sure the block is below content.height to stack components in the box
	box.dimensions.content.y = containingBlock.content.height + containingBlock.content.y +
		box.dimensions.margin.top + box.dimensions.border.top + box.dimensions.padding.top
}

func (box *LayoutBox) layoutBlockChildren() {
	for _, child := range box.children {
		child.Layout(box.dimensions)
		// Track the height so each child is laid out below the previous content
		box.dimensions.content.height = box.dimensions.content.height + child.dimensions.marginBox().height
	}
}

func (box *LayoutBox) calculateHeight() {
	// If the height is set to an explicit length, use that exact length
	// Otherwise, just keep the value set by layoutBlockChildren()
	if height, ok := box.styledNode.Value("height"); ok {
		if height.Length.Unit == css.Px {
			box.dimensions.content.height = height.Length.Quantity
		}
	}
}
