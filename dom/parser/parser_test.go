package parser

import (
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

func TestParseElement(t *testing.T) {
	input := `<a></a><b></b><c></c>`
	tests := []*dom.Node{
		{Type: dom.NodeElement, Tag: "a"},
		{Type: dom.NodeElement, Tag: "b"},
		{Type: dom.NodeElement, Tag: "c"},
	}

	l := lexer.New(input)
	p := New(l)

	nodes := p.Parse()
	t.Log(p.Errors())
	if nodes == nil {
		t.Fatal("fail to parse DOM")
	}
	if len(nodes) != len(tests) {
		t.Fatalf("invalide node count. expected=%d got=%d", len(tests), len(nodes))
	}
	for i, tt := range tests {
		n := nodes[i]
		if n.Tag != tt.Tag {
			t.Fatalf("tests[%d]: invalid tag. expected=%q got=%q", i, tt.Tag, n.Tag)
		}
		if n.Type != tt.Type {
			t.Fatalf("tests[%d]: invalid type. expected=%q got=%q", i, tt.Type, n.Type)
		}
	}
}

func TestParseElement_recursion(t *testing.T) {
	d := &dom.Node{Type: dom.NodeElement, Tag: "d"}
	c := &dom.Node{Type: dom.NodeElement, Tag: "c", NextSibling: d}
	b := &dom.Node{Type: dom.NodeElement, Tag: "b", NextSibling: c}
	a := &dom.Node{Type: dom.NodeElement, Tag: "a", FirstChild: b, LastChild: d}
	d.PrevSibling = c
	c.PrevSibling = b
	b.Parent, c.Parent, d.Parent = a, a, a

	input := `<a><b></b><c></c><d></d></a>`
	tests := []*dom.Node{a}

	l := lexer.New(input)
	p := New(l)

	nodes := p.Parse()
	t.Log(p.Errors())
	if nodes == nil {
		t.Fatal("fail to parse DOM")
	}
	if len(nodes) != len(tests) {
		t.Fatalf("invalide node count. expected=%d got=%d", len(tests), len(nodes))
	}

	tt := tests[0].FirstChild
	n := nodes[0].FirstChild
	for tt != nil {
		if n == nil {
			t.Fatalf("tests[0][%q]: missing node.", tt.Tag)
		}
		if n.Tag != tt.Tag {
			t.Fatalf("tests[0][%q]: invalid tag. expected=%q got=%q", tt.Tag, tt.Tag, n.Tag)
		}
		if n.Type != tt.Type {
			t.Fatalf("tests[0][%q]: invalid type. expected=%q got=%q", tt.Tag, tt.Type, n.Type)
		}

		tt = tt.NextSibling
		n = n.NextSibling
	}
}
