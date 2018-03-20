// Package layout represents the layout tree of the HTML renderer.
//
// The main type is called LayoutBox.
package layout

import (
	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/style"
)

const (
	BlockNode BoxType = iota
	InlineNode
	AnonymousBlock
)

type BoxType int

// LayoutBox is the building block of the layout tree, associated to one StyleNode
type LayoutBox struct {
	// Dimensions is the box position, size, margin, padding and border
	Dimensions Dimensions

	// BoxType is the type of the box (inline, block, anonymous)
	BoxType BoxType

	// StyleNode holds the style node wrapped by this layout node
	StyledNode *style.StyledNode

	// Children of this node, in the layout tree, following the structure of the style tree
	Children []*LayoutBox
}

// Dimensions represents the position, size, margin, padding and border of a layout box
type Dimensions struct {
	// Content represents the position and size of the content area relative to the document origin:
	Content Rect

	// Surrounding edges:
	padding, Border, margin EdgeSizes
}

// Rect represents the position and size of a box on the screen
type Rect struct {
	X, Y, Width, Height float64
}

// EdgeSizes is a placeholder for four float values
type EdgeSizes struct {
	Left, Right, Top, Bottom float64
}

// marginBox returns the area covered by the content area plus padding, borders, and margin
func (d Dimensions) marginBox() Rect {
	return d.BorderBox().expandedBy(d.margin)
}

// BorderBox returns the area covered by the content area plus padding and borders
func (d Dimensions) BorderBox() Rect {
	return d.paddingBox().expandedBy(d.Border)
}

// paddingBox returns the area covered by the content area plus its padding
func (d Dimensions) paddingBox() Rect {
	return d.Content.expandedBy(d.padding)
}

func (r Rect) expandedBy(edge EdgeSizes) Rect {
	return Rect{
		X:      r.X - edge.Left,
		Y:      r.Y - edge.Top,
		Width:  r.Width + edge.Left + edge.Right,
		Height: r.Height + edge.Top + edge.Bottom,
	}
}

func newLayoutBox(boxType BoxType, styledNode *style.StyledNode) *LayoutBox {
	var children []*LayoutBox
	return &LayoutBox{
		BoxType:    boxType,
		StyledNode: styledNode,
		Children:   children,
	}
}

func GenerateLayoutTree(styleTree *style.StyledNode) *LayoutBox {
	var boxType BoxType
	switch styleTree.Display() {
	case style.Inline:
		boxType = InlineNode
	case style.Block:
		boxType = BlockNode
	case style.None:
		panic("Root StyledNode has display:none")
	}

	root := newLayoutBox(boxType, styleTree)

	for _, child := range styleTree.Children {
		switch child.Display() {
		case style.Inline:
			ic := root.getInlineContainer()
			ic.Children = append(ic.Children, GenerateLayoutTree(child))
		case style.Block:
			root.Children = append(root.Children, GenerateLayoutTree(child))
		case style.None:
			// Skip
		}
	}

	return root
}

// If a block node contains an inline child:
func (box *LayoutBox) getInlineContainer() *LayoutBox {
	switch box.BoxType {
	case InlineNode:
		fallthrough
	case AnonymousBlock:
		return box
	case BlockNode:
		// If we've just generated an anonymous block box, keep using it.
		// Otherwise, create a new one.
		if len(box.Children) == 0 {
			box.Children = append(box.Children, newLayoutBox(AnonymousBlock, nil))
			return box.Children[0]
		}

		lastChild := box.Children[len(box.Children)-1]
		switch lastChild.BoxType {
		case AnonymousBlock:
			return lastChild
		default:
			box.Children = append(box.Children, newLayoutBox(AnonymousBlock, nil))
			return box.Children[len(box.Children)-1]
		}
	}
	panic("No more cases to switch")
}

func (box *LayoutBox) Layout(containingBlock Dimensions) {
	switch box.BoxType {
	case InlineNode:
		// TODO
		panic("Inline Node Unimplemented")
	case BlockNode:
		box.layoutBlock(containingBlock)
	case AnonymousBlock:
		// TODO
		panic("Anonymous Block Unimplemented")
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
	style := box.StyledNode

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
	if width != auto && total > containingBlock.Content.Width {
		if marginLeft == auto {
			marginLeft = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}
		if marginRight == auto {
			marginRight = css.Value{Length: css.Length{Quantity: 0.0, Unit: css.Px}}
		}
	}

	// Check for over or underflow, and adjust "auto" dimensions accordingly
	underflow := containingBlock.Content.Width - total

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

	box.Dimensions.Content.Width = width.ToPx()

	box.Dimensions.padding.Left = paddingLeft.ToPx()
	box.Dimensions.padding.Right = paddingRight.ToPx()

	box.Dimensions.Border.Left = borderLeft.ToPx()
	box.Dimensions.Border.Right = borderRight.ToPx()

	box.Dimensions.margin.Left = marginLeft.ToPx()
	box.Dimensions.margin.Right = marginRight.ToPx()
}

func (box *LayoutBox) calculatePosition(containingBlock Dimensions) {
	style := box.StyledNode

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

	box.Dimensions.margin.Top = marginTop.ToPx()
	box.Dimensions.margin.Bottom = marginBottom.ToPx()
	box.Dimensions.Border.Top = borderTop.ToPx()
	box.Dimensions.Border.Bottom = borderBottom.ToPx()
	box.Dimensions.padding.Top = paddingTop.ToPx()
	box.Dimensions.padding.Bottom = paddingBottom.ToPx()

	box.Dimensions.Content.X = containingBlock.Content.X +
		box.Dimensions.margin.Left + box.Dimensions.Border.Left + box.Dimensions.padding.Left

	// Position the box below all the previous boxes in the container.
	// Making sure the block is below content.height to stack components in the box
	box.Dimensions.Content.Y = containingBlock.Content.Height + containingBlock.Content.Y +
		box.Dimensions.margin.Top + box.Dimensions.Border.Top + box.Dimensions.padding.Top
}

func (box *LayoutBox) layoutBlockChildren() {
	for _, child := range box.Children {
		child.Layout(box.Dimensions)
		// Track the height so each child is laid out below the previous content
		box.Dimensions.Content.Height = box.Dimensions.Content.Height + child.Dimensions.marginBox().Height
	}
}

func (box *LayoutBox) calculateHeight() {
	// If the height is set to an explicit length, use that exact length
	// Otherwise, just keep the value set by layoutBlockChildren()
	if height, ok := box.StyledNode.Value("height"); ok {
		if height.Length.Unit == css.Px {
			box.Dimensions.Content.Height = height.Length.Quantity
		}
	}
}
