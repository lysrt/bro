package main

import (
	"strings"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/dom"
	"golang.org/x/net/html"
)

// PropertyMap maps CSS properties to Value.
type PropertyMap map[string]css.Value

// StyledNode represents a DOM Node with the associated CSS.
type StyledNode struct {
	Node            *html.Node
	SpecifiedValues PropertyMap
	Children        []StyledNode
}

// Value returns the value of a CSS property if it exists
func (node *StyledNode) Value(property string) (value css.Value, ok bool) {
	value, ok = node.SpecifiedValues[property]
	return
}

type Display string

const (
	Inline Display = "inline"
	Block  Display = "block"
	None   Display = "None"
)

// Display returns the type of CSS display of a StyledNode
func (node *StyledNode) Display() Display {
	value, ok := node.Value("display")
	if !ok {
		return Inline
	}

	switch strings.ToLower(value.Keyword) {
	case "block":
		return Block
	case "none":
		return None
	default:
		return Inline
	}
}

// MatchedRule represents a matched rule with a given specificity.
type MatchedRule struct {
	Rule        css.Rule
	Specificity css.Specificity
}

func GenerateStyleTree(root *html.Node, css *css.Stylesheet) *StyledNode {
	var propertyMap PropertyMap

	switch root.Type {
	case html.ElementNode:
		if css == nil {
			propertyMap = make(PropertyMap)
		} else {
			propertyMap = specifiedValues(root, css)
		}
	case html.TextNode:
		propertyMap = make(PropertyMap)
	}

	var children []StyledNode
	for _, child := range dom.NodeChildren(root) {
		styled := GenerateStyleTree(child, css)
		children = append(children, *styled)
	}

	return &StyledNode{
		Node:            root,
		SpecifiedValues: propertyMap,
		Children:        children,
	}
}

// specifiedValues returns the apply properties of a node.
func specifiedValues(element *html.Node, stylesheet *css.Stylesheet) PropertyMap {
	properties := make(PropertyMap)
	rules := matchingRules(element, stylesheet)

	// Order from lowest to highest specificity
	for i := range rules {
		for j := range rules[i:] {
			speI := rules[i].Specificity
			speJ := rules[j].Specificity
			if speI.A > speJ.A {
				continue
			}
			if speI.A < speJ.A {
				rules[i], rules[j] = rules[j], rules[i]
				continue
			}
			if speI.B > speJ.B {
				continue
			}
			if speI.B < speJ.B {
				rules[i], rules[j] = rules[j], rules[i]
				continue
			}
			if speI.C > speJ.C {
				continue
			}
			if speI.C < speJ.C {
				rules[i], rules[j] = rules[j], rules[i]
				continue
			}
		}
	}

	for _, r := range rules {
		for _, d := range r.Rule.Declarations {
			properties[d.Name] = d.Value
		}
	}

	return properties
}

// matchingRules returns all the matched rules for a node.
func matchingRules(n *html.Node, stylesheet *css.Stylesheet) []MatchedRule {
	var matches []MatchedRule
	for _, r := range stylesheet.Rules {
		m, ok := matchRule(n, r)
		if !ok {
			continue
		}
		matches = append(matches, m)
	}
	return matches
}

// matchRule tries to match a rule to a node and return the most specific one.
func matchRule(n *html.Node, rule css.Rule) (m MatchedRule, ok bool) {
	for _, s := range rule.Selectors {
		if matchSelector(n, s) {
			ok = true
			m = MatchedRule{
				Rule:        rule,
				Specificity: s.Specificity(),
			}
			return
		}
	}
	return
}

// matchSelector matches a node with a selector.
// Their is a match only if all the field of the selector find a match.
func matchSelector(n *html.Node, selector css.Selector) bool {
	if selector.TagName != "" && n.Data != selector.TagName {
		return false
	}
	if selector.ID != "" && selector.ID != dom.NodeGetID(n) {
		return false
	}
outer_loop:
	for _, c := range selector.Classes {
		for _, cc := range dom.NodeGetClasses(n) {
			if c == cc {
				continue outer_loop
			}
		}
		return false
	}
	return true
}
