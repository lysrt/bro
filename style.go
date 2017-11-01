package main

import "golang.org/x/net/html"

// PropertyMap maps CSS properties to Value.
type PropertyMap map[string]Value

// StyledNode represents a DOM Node with the associated CSS.
type StyledNode struct {
	Node            *html.Node
	SpecifiedValues PropertyMap
	Children        []StyledNode
}

func matches(n *html.Node, selector *Selector) bool {
	if selector.TagName != "" {
		if n.Data == selector.TagName {
			return true
		}
		return false
	}
	if selector.ID != "" {
		if selector.ID == NodeGetID(n) {
			return true
		}
		return false
	}
	if selector.Classes != nil {
		for _, c := range NodeGetClasses(n) {
			if selector.Classes[0] == c {
				return true
			}
		}
	}
	return false
}
