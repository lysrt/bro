package parser

import (
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

func isNodeEqual(a, b *dom.Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a.Type != b.Type {
		return false
	}
	if a.Tag != b.Tag {
		return false
	}
	if (a.Parent == nil || b.Parent == nil) && a.Parent != b.Parent {
		return false
	}
	if (a.FirstChild == nil || b.FirstChild == nil) && a.FirstChild != b.FirstChild {
		return false
	}
	if (a.LastChild == nil || b.LastChild == nil) && a.LastChild != b.LastChild {
		return false
	}
	if (a.NextSibling == nil || b.NextSibling == nil) && a.NextSibling != b.NextSibling {
		return false
	}
	if (a.PrevSibling == nil || b.PrevSibling == nil) && a.PrevSibling != b.PrevSibling {
		return false
	}
	return true
}

func walkNodes(a, b *dom.Node, fn func(a, b *dom.Node)) {
	if a == nil || b == nil {
		return
	}
	fn(a, b)
	walkNodes(a.FirstChild, b.FirstChild, fn)
	walkNodes(a.NextSibling, b.NextSibling, fn)
}

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

func TestParseElement_recurse(t *testing.T) {
	f := &dom.Node{Type: dom.NodeElement, Tag: "f"}
	e := &dom.Node{Type: dom.NodeElement, Tag: "e", NextSibling: f}
	d := &dom.Node{Type: dom.NodeElement, Tag: "d"}
	c := &dom.Node{Type: dom.NodeElement, Tag: "c", NextSibling: d}
	b := &dom.Node{Type: dom.NodeElement, Tag: "b", NextSibling: c, FirstChild: e, LastChild: f}
	a := &dom.Node{Type: dom.NodeElement, Tag: "a", FirstChild: b, LastChild: d}
	f.PrevSibling = e
	d.PrevSibling = c
	c.PrevSibling = b
	e.Parent, f.Parent = b, b
	b.Parent, c.Parent, d.Parent = a, a, a

	input := `<a><b><e></e><f></f></b><c></c><d></d></a>`
	tests := []*dom.Node{a}

	l := lexer.New(input)
	p := New(l)

	nodes := p.Parse()
	t.Logf("parser errors: %v", p.Errors())
	if nodes == nil {
		t.Fatal("fail to parse DOM")
	}
	if len(nodes) != len(tests) {
		t.Fatalf("invalide node count. expected=%d got=%d", len(tests), len(nodes))
	}

	walkNodes(tests[0], nodes[0], func(a, b *dom.Node) {
		if !isNodeEqual(a, b) {
			t.Logf("a=%v", a)
			t.Logf("b=%v", b)
			t.Fatal("nodes are different")
		}
	})
}

func TestParseElement_attributes(t *testing.T) {
	input := `<a class="awesome"></a><b id="unique" class="awesome"></b><c id="intimidating" class="awesome"></c>`
	tests := []*dom.Node{
		{
			Type: dom.NodeElement,
			Tag:  "a",
			Attributes: map[string]string{
				"class": "awesome",
			},
		},
		{
			Type: dom.NodeElement,
			Tag:  "b",
			Attributes: map[string]string{
				"class": "awesome",
				"id":    "unique",
			},
		},
		{
			Type: dom.NodeElement,
			Tag:  "c",
			Attributes: map[string]string{
				"class": "awesome",
				"id":    "intimidating",
			},
		},
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
		for k, v := range tt.Attributes {
			vv, ok := n.Attributes[k]
			if !ok {
				t.Fatalf("tests[%d]: missing attribute %q.", i, k)
			}
			if v != vv {
				t.Fatalf("tests[%d]: bad attribute %q. expected=%q got=%q", i, k, vv, v)
			}
		}
	}
}

func TestParseText(t *testing.T) {
	input := `<a>I can read text!</a>`
	tests := *dom.Node{
		{Type: dom.NodeElement, Tag: "a"},
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
