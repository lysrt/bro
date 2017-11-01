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

func specifiedValues(n *html.Node, r Rule) (m MatchedRule, ok bool) {
	var matches []MatchedRule
	for _, s := range r.Selectors {
		if !matchSelector(n, s) {
			continue
		}
		matches = append(matches, MatchedRule{
			Rule:        r,
			Specificity: s.Specificity(),
		})
	}
	if len(matches) == 0 {
		return
	}
	for _, match := range matches {
		if match.Specificity.A < m.Specificity.A {
			continue
		}
		if match.Specificity.A > m.Specificity.A {
			m = match
			continue
		}
		if match.Specificity.B < m.Specificity.B {
			continue
		}
		if match.Specificity.B > m.Specificity.B {
			m = match
			continue
		}
		if match.Specificity.C < m.Specificity.C {
			continue
		}
		if match.Specificity.C > m.Specificity.C {
			m = match
			continue
		}
	}
	ok = true
	return
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
