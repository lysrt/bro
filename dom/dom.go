package dom

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// NodeGetID extracts ID from a node.
func NodeGetID(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key != "id" {
			continue
		}
		return attr.Val
	}
	return ""
}

// NodeGetClasses extracts classes from a node.
func NodeGetClasses(n *html.Node) []string {
	for _, attr := range n.Attr {
		if attr.Key != "class" {
			continue
		}
		return strings.Fields(attr.Val)
	}
	return nil
}

// NodeChildren returns all child nodes of n
func NodeChildren(n *html.Node) []*html.Node {
	var children []*html.Node

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
func NodeFirstElementChild(n *html.Node) *html.Node {
	for e := n.FirstChild; e != nil; e = e.NextSibling {
		if e.Type == html.TextNode {
			continue
		}
		return e
	}
	return nil
}

// NodeLastElementChild returns the last element child of the node.
func NodeLastElementChild(n *html.Node) *html.Node {
	for e := n.LastChild; e != nil; e = e.PrevSibling {
		if e.Type == html.TextNode {
			continue
		}
		return e
	}
	return nil
}

// NodeNextElementSibling returns the next element sibling of the node.
func NodeNextElementSibling(n *html.Node) *html.Node {
	for e := n.NextSibling; e != nil; e = e.NextSibling {
		if e.Type == html.TextNode {
			continue
		}
		return e
	}
	return nil
}

// NodePrevElementSibling returns the previous element sibling of the node.
func NodePrevElementSibling(n *html.Node) *html.Node {
	for e := n.PrevSibling; e != nil; e = e.PrevSibling {
		if e.Type == html.TextNode {
			continue
		}
		return e
	}
	return nil
}

func Parcour(n *html.Node) {
	ParcourN(n, 0)
}

func ParcourN(n *html.Node, depth int) {
	current := n

	if current.Type == html.TextNode {
		text := current.Data
		text = strings.TrimSpace(text)
		if text != "" {
			p(depth, "-> %v \"%v\"\n", printNodeType(current.Type), text)
		}
	} else if current.Type == html.ElementNode {
		if len(current.Attr) > 0 {
			p(depth, "-> %v <%v> %v\n", printNodeType(current.Type), current.DataAtom, current.Attr)
		} else {
			p(depth, "-> %v <%v>\n", printNodeType(current.Type), current.DataAtom)
		}
		if current.Data != current.DataAtom.String() {
			p(depth, "Tag: %v\n", current.Data) // equals atom, if atom is recognized
		}
	} else {
		p(depth, "-> %v \"%v\"\n", printNodeType(current.Type), current.Data)
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

func printNodeType(t html.NodeType) string {
	switch t {
	case html.DoctypeNode:
		return "Doctype"
	case html.CommentNode:
		return "Comment"
	case html.DocumentNode:
		return "Document"
	case html.ElementNode:
		return "Element"
	case html.ErrorNode:
		return "Error"
	case html.TextNode:
		return "Text"
	default:
		return "UNKNOWN"
	}
}
