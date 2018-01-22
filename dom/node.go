package dom

type Node interface {
	Parent() Node
	FirstChild() Node
	LastChild() Node
	PrevSibling() Node
	NextSibling() Node
}

type baseNode struct {
	parent      Node
	firstChild  Node
	lastChild   Node
	prevSibling Node
	nextSibling Node
}

func (n *baseNode) Parent() Node      { return n.parent }
func (n *baseNode) FirstChild() Node  { return n.firstChild }
func (n *baseNode) LastChild() Node   { return n.lastChild }
func (n *baseNode) PrevSibling() Node { return n.prevSibling }
func (n *baseNode) NextSibling() Node { return n.nextSibling }

func (n *baseNode) SetParent(nn Node)      { n.parent = nn }
func (n *baseNode) SetFirstChild(nn Node)  { n.firstChild = nn }
func (n *baseNode) SetLastChild(nn Node)   { n.lastChild = nn }
func (n *baseNode) SetPrevSibling(nn Node) { n.prevSibling = nn }
func (n *baseNode) SetNextSibling(nn Node) { n.nextSibling = nn }

type Attribute struct {
	Key string
	Val string
}

type Document struct {
	baseNode
	Children []Node
}

type Text struct {
	baseNode
}

type Element struct {
	baseNode
	Tag        string
	Attributes []Attribute
}
