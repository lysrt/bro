package parser

import (
	"testing"

	"github.com/lysrt/bro/html"
	"github.com/lysrt/bro/html/lexer"
)

func isNodeEqual(a, b *html.Node) bool {
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

// compareNodes walks throught a and b calling fn on each iteration.
// The function returns fall if a & b does not have the same number of element.
func compareNodes(a, b *html.Node, fn func(a, b *html.Node)) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	fn(a, b)

	if !compareNodes(a.FirstChild, b.FirstChild, fn) {
		return false
	}
	if !compareNodes(a.NextSibling, b.NextSibling, fn) {
		return false
	}
	return true
}

func TestParseElement(t *testing.T) {
	expected := &html.Node{Type: html.NodeElement, Tag: "html"}
	{
		head := &html.Node{Type: html.NodeElement, Tag: "head"}
		body := &html.Node{Type: html.NodeElement, Tag: "body"}
		body.AddChild(&html.Node{Type: html.NodeElement, Tag: "a"})
		body.AddChild(&html.Node{Type: html.NodeElement, Tag: "b"})
		body.AddChild(&html.Node{Type: html.NodeElement, Tag: "c"})
		expected.AddChild(head)
		expected.AddChild(body)
	}

	input := `<a></a><b></b><c></c>`
	l := lexer.New(input)
	p := New(l)

	parsed := p.Parse()
	t.Log(p.Errors())
	if parsed == nil {
		t.Fatal("fail to parse DOM")
	}
	ok := compareNodes(expected, parsed, func(a, b *html.Node) {
		if a.Tag != b.Tag {
			t.Fatalf("invalid tag. expected=%q got=%q", a.Tag, b.Tag)
		}
		if a.Type != b.Type {
			t.Fatalf("invalid type. expected=%q got=%q", a.Type, b.Type)
		}
	})
	if !ok {
		t.Fatal("fail to compare nodes")
	}
}

func TestParseElement_recurse(t *testing.T) {
	a := &html.Node{Type: html.NodeElement, Tag: "a"}
	{
		b := &html.Node{Type: html.NodeElement, Tag: "b"}
		b.AddChild(&html.Node{Type: html.NodeElement, Tag: "e"})
		b.AddChild(&html.Node{Type: html.NodeElement, Tag: "f"})
		a.AddChild(b)
		a.AddChild(&html.Node{Type: html.NodeElement, Tag: "c"})
		a.AddChild(&html.Node{Type: html.NodeElement, Tag: "d"})
	}
	expected := &html.Node{Type: html.NodeElement, Tag: "html"}
	{
		head := &html.Node{Type: html.NodeElement, Tag: "head"}
		body := &html.Node{Type: html.NodeElement, Tag: "body"}
		body.AddChild(a)
		expected.AddChild(head)
		expected.AddChild(body)
	}

	input := `<a><b><e></e><f></f></b><c></c><d></d></a>`
	l := lexer.New(input)
	p := New(l)

	parsed := p.Parse()
	t.Log(p.Errors())
	if parsed == nil {
		t.Fatal("fail to parse DOM")
	}

	ok := compareNodes(expected, parsed, func(a, b *html.Node) {
		if !isNodeEqual(a, b) {
			t.Logf("a=%v", a)
			t.Logf("b=%v", b)
			t.Fatal("nodes are different")
		}
	})
	if !ok {
		t.Fatal("fail to compare nodes")
	}
}

func TestParseElement_attributes(t *testing.T) {
	expected := &html.Node{Type: html.NodeElement, Tag: "html"}
	{
		expected.AddChild(&html.Node{Type: html.NodeElement, Tag: "head"})
		body := &html.Node{Type: html.NodeElement, Tag: "body"}
		body.AddChild(&html.Node{
			Type: html.NodeElement,
			Tag:  "a",
			Attributes: map[string]string{
				"class": "cool",
			},
		})
		body.AddChild(&html.Node{
			Type: html.NodeElement,
			Tag:  "b",
			Attributes: map[string]string{
				"id": "unique",
			},
		})
		body.AddChild(&html.Node{
			Type: html.NodeElement,
			Tag:  "c",
			Attributes: map[string]string{
				"id":    "crazy",
				"class": "cool",
			},
		})
		expected.AddChild(body)
	}

	input := `<a class="cool"></a><b id="unique"></b><c id="crazy" class="cool"></c>`
	l := lexer.New(input)
	p := New(l)

	parsed := p.Parse()
	t.Log(p.Errors())
	if parsed == nil {
		t.Fatal("fail to parse DOM")
	}
	ok := compareNodes(expected, parsed, func(a, b *html.Node) {
		if b.Tag != a.Tag {
			t.Fatalf("invalid tag. expected=%q got=%q", a.Tag, b.Tag)
		}
		if b.Type != a.Type {
			t.Fatalf("invalid type. expected=%q got=%q", a.Type, b.Type)
		}
		for k, v := range a.Attributes {
			vv, ok := b.Attributes[k]
			if !ok {
				t.Fatalf("missing attribute %q.", k)
			}
			if v != vv {
				t.Fatalf("bad attribute %q. expected=%q got=%q", k, vv, v)
			}
		}
	})
	if !ok {
		t.Fatal("fail to compare nodes")
	}
}

func TestParseText(t *testing.T) {
	expected := &html.Node{Type: html.NodeElement, Tag: "html"}
	{
		expected.AddChild(&html.Node{Type: html.NodeElement, Tag: "head"})
		body := &html.Node{Type: html.NodeElement, Tag: "body"}
		a := &html.Node{Type: html.NodeElement, Tag: "a"}
		a.AddChild(&html.Node{Type: html.NodeText, TextContent: "I can read text!"})
		body.AddChild(a)
		expected.AddChild(body)
	}

	input := `<a>I can read text!</a>`
	l := lexer.New(input)
	p := New(l)

	parsed := p.Parse()
	t.Log(p.Errors())
	if parsed == nil {
		t.Fatal("fail to parse DOM")
	}
	ok := compareNodes(expected, parsed, func(a, b *html.Node) {
		if b.Type != a.Type {
			t.Fatalf("invalid type. expected=%q got=%q", a.Type, b.Type)
		}
		if b.TextContent != a.TextContent {
			t.Fatalf("invalid text content. expected=%q got=%q", a.TextContent, b.TextContent)
		}
	})
	if !ok {
		t.Fatal("fail to compare nodes")
	}
}
