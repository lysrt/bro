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
	if selector.TagName != "" && n.Data != selector.TagName {
		return false
	}
	if selector.ID != "" && selector.ID != NodeGetID(n) {
		return false
	}
outer_loop:
	for _, c := range selector.Class {
		for _, cc := range NodeGetClasses(n) {
			if c == cc {
				continue outer_loop
			}
		}
		return false
	}
	return true
}
