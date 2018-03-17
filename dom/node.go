package dom

import "fmt"

const (
	NodeError   NodeType = "error"
	NodeText    NodeType = "text"
	NodeElement NodeType = "element"
	NodeComment NodeType = "comment"
)

type NodeType string

type Node struct {
	Parent      *Node
	FirstChild  *Node
	LastChild   *Node
	PrevSibling *Node
	NextSibling *Node

	Type NodeType

	Tag        string
	Attributes map[string]string

	TextContent string
}

func (n *Node) String() string {
	return fmt.Sprintf(
		"{Tag:%q, Parent:%p, FirstChild:%p, LastChild:%p, PrevSibling:%p, NextSibling:%p}",
		n.Tag,
		n.Parent,
		n.FirstChild,
		n.LastChild,
		n.PrevSibling,
		n.NextSibling,
	)
}
