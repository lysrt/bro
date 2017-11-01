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

// MatchedRule represents a matched rule with a given specificity.
type MatchedRule struct {
	Rule        Rule
	Specificity Specificity
}

// matchRule tries to match a rule to a node and return the most specifique one.
func matchRule(n *html.Node, r Rule) (m MatchedRule, ok bool) {
	for _, s := range r.Selectors {
		if matchSelector(n, s) {
			ok = true
			m = MatchedRule{
				Rule:        r,
				Specificity: s.Specificity(),
			}
			return
		}
	}
	return
}

// matchingRules returns all the matched rules for a node.
func matchingRules(n *html.Node, stylesheet Stylesheet) []MatchedRule {
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

// specifiedValues returns the apply properties of a node.
func specifiedValues(n *html.Node, stylesheet Stylesheet) PropertyMap {
	properties := make(PropertyMap)
	rules := matchingRules(n, stylesheet)

	// order from lowest to highest
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
		for _, d := range r.Rule.Declatations {
			properties[d.name] = d.value
		}
	}

	return properties
}

// matchSelector matches a node with a selector.
// Their is a match only if all the field of the selector find a match.
func matchSelector(n *html.Node, selector Selector) bool {
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
