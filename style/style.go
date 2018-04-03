package style

import (
	"strings"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/html"
)

// PropertyMap is used by a node to keep track of its applied CSS declarations.
// For instance, the node <body> could have the corresponding map:
// {
//     "margin-left": css.Value{Length: css.Length{Quantity: 5, Unit: css.Px}},
//     "color": css.Value{Color: css.Color{255, 255, 0, 0}}
// }
type PropertyMap map[string]css.Value

// StyledNode represents a DOM Node with the associated CSS.
type StyledNode struct {
	Node            *html.Node
	SpecifiedValues PropertyMap
	Children        []*StyledNode
}

// Value returns the value of a given CSS property name of a StyleNode, if it has one
func (node *StyledNode) Value(property string) (value css.Value, ok bool) {
	value, ok = node.SpecifiedValues[property]
	return
}

type Display string

const (
	Inline Display = "inline"
	Block  Display = "block"
	None   Display = "none"
)

// Display returns the CSS display type of a StyledNode
func (node *StyledNode) Display() Display {
	value, ok := node.Value("display")
	if !ok {
		return Block // Block is the default display type
	}

	switch strings.ToLower(value.Keyword) {
	case "block":
		return Block
	case "none":
		return None
	default:
		return Inline // Any value other than block or none will result in inline display
	}
}

// GenerateStyleTree a DOM node and its children with CSS rules from a Stylesheet.
func GenerateStyleTree(root *html.Node, css *css.Stylesheet) *StyledNode {
	var propertyMap PropertyMap

	switch root.Type {
	case html.NodeElement:
		if css == nil {
			propertyMap = make(PropertyMap)
		} else {
			propertyMap = specifiedValues(root, css)
		}
	case html.NodeText:
		propertyMap = make(PropertyMap)
	}

	var children []*StyledNode
	for _, child := range html.NodeChildren(root) {
		styled := GenerateStyleTree(child, css)
		children = append(children, styled)
	}

	return &StyledNode{
		Node:            root,
		SpecifiedValues: propertyMap,
		Children:        children,
	}
}

// MatchedRule represents a matched rule with a given specificity.
type MatchedRule struct {
	Rule        css.Rule
	Specificity css.Specificity
}

// specifiedValues returns a map of all CSS properties applied to a given DOM node.
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

	// If several rules have the same name, highest specificity rules will override low specificity ones
	for _, r := range rules {
		for _, d := range r.Rule.Declarations {
			properties[d.Name] = d.Value
		}
	}

	return properties
}

// matchingRules returns all the matched CSS rules for a given DOM node.
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

// matchRule tries to match a CSS rule to a DOM node and returns the most specific one.
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

// matchSelector tries to match a DOM node with a CSS selector.
// There is a match only if all the fields of the selector match.
func matchSelector(n *html.Node, selector css.Selector) bool {
	if selector.TagName != "" && selector.TagName == "*" {
		return true
	}
	if selector.TagName != "" && n.Tag != selector.TagName {
		return false
	}
	if selector.ID != "" && selector.ID != html.NodeGetID(n) {
		return false
	}
outer_loop:
	for _, c := range selector.Classes {
		for _, cc := range html.NodeGetClasses(n) {
			if c == cc {
				continue outer_loop
			}
		}
		// The CSS selector class is not one of the DOM node classes
		return false
	}
	return true
}
