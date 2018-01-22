package dom

type Node interface {
	Parent() Node
	FirstChild() Node
	LastChild() Node
	PrevSibling() Node
	NextSibling() Node
}

type Attribute struct {
	Key string
	Val string
}

type HTMLNode struct {
	parent      Node
	firstChild  Node
	lastChild   Node
	prevSibling Node
	nextSibling Node
}

func (n *HTMLNode) Parent() Node      { return n.parent }
func (n *HTMLNode) FirstChild() Node  { return n.firstChild }
func (n *HTMLNode) LastChild() Node   { return n.lastChild }
func (n *HTMLNode) PrevSibling() Node { return n.prevSibling }
func (n *HTMLNode) NextSibling() Node { return n.nextSibling }

func (n *HTMLNode) SetParent(nn Node)      { n.parent = nn }
func (n *HTMLNode) SetFirstChild(nn Node)  { n.firstChild = nn }
func (n *HTMLNode) SetLastChild(nn Node)   { n.lastChild = nn }
func (n *HTMLNode) SetPrevSibling(nn Node) { n.prevSibling = nn }
func (n *HTMLNode) SetNextSibling(nn Node) { n.nextSibling = nn }

type TextNode struct {
	HTMLNode
}

type ElementNode struct {
	HTMLNode
	Attributes []Attribute
}
