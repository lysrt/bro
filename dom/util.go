package dom

import (
	"fmt"
	"strings"
)

// NodeGetID extracts ID from a node.
func NodeGetID(n *Node) string {
	for k, v := range n.Attributes {
		if k != "id" {
			continue
		}
		return v
	}
	return ""
}

// NodeGetClasses extracts classes from a node.
func NodeGetClasses(n *Node) []string {
	for k, v := range n.Attributes {
		if k != "class" {
			continue
		}
		return strings.Fields(v)
	}
	return nil
}

// NodeChildren returns all child nodes of n
func NodeChildren(n *Node) []*Node {
	var children []*Node

	f := n.FirstChild
	if f == nil {
		return children
	}

	children = append(children, f)

	next := f.NextSibling
	for next != nil {
		children = append(children, next)
		next = next.NextSibling
	}

	return children
}

// NodeFirstElementChild returns the first element child of the node.
func NodeFirstElementChild(n *Node) *Node {
	for e := n.FirstChild; e != nil; e = e.NextSibling {
		if e.Type == NodeText {
			continue
		}
		return e
	}
	return nil
}

// NodeLastElementChild returns the last element child of the node.
func NodeLastElementChild(n *Node) *Node {
	for e := n.LastChild; e != nil; e = e.PrevSibling {
		if e.Type == NodeText {
			continue
		}
		return e
	}
	return nil
}

// NodeNextElementSibling returns the next element sibling of the node.
func NodeNextElementSibling(n *Node) *Node {
	for e := n.NextSibling; e != nil; e = e.NextSibling {
		if e.Type == NodeText {
			continue
		}
		return e
	}
	return nil
}

// NodePrevElementSibling returns the previous element sibling of the node.
func NodePrevElementSibling(n *Node) *Node {
	for e := n.PrevSibling; e != nil; e = e.PrevSibling {
		if e.Type == NodeText {
			continue
		}
		return e
	}
	return nil
}

func Parcour(n *Node) {
	ParcourN(n, 0)
}

func ParcourN(n *Node, depth int) {
	current := n

	if current.Type == NodeText {
		text := current.TextContent
		text = strings.TrimSpace(text)
		if text != "" {
			p(depth, "-> %v \"%v\"\n", printNodeType(current.Type), text)
		}
	} else if current.Type == NodeElement {
		if len(current.Attributes) > 0 {
			p(depth, "-> %v <%v> %v\n", printNodeType(current.Type), current.Tag, current.Attributes)
		} else {
			p(depth, "-> %v <%v>\n", printNodeType(current.Type), current.Tag)
		}
		if current.TextContent != current.Tag {
			p(depth, "Tag: %v\n", current.TextContent) // equals atom, if atom is recognized
		}
	} else {
		p(depth, "-> %v \"%v\"\n", printNodeType(current.Type), current.TextContent)
	}

	if current.FirstChild != nil {
		ParcourN(current.FirstChild, depth+1)
	}

	if current.NextSibling != nil {
		ParcourN(current.NextSibling, depth)
	}
}

func p(depth int, format string, args ...interface{}) {
	var spaces string
	for i := 0; i < depth; i++ {
		spaces += "   "
	}
	fmt.Printf(spaces+format, args...)
}

func printNodeType(t NodeType) string {
	return string(t)
}
