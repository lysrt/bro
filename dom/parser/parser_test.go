package parser

import (
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

func TestDOMElement(t *testing.T) {
	input := `<a></a><b></b><c></c>`
	tests := []*dom.Node{
		{Tag: "a"},
		{Tag: "b"},
		{Tag: "c"},
	}

	l := lexer.New(input)
	p := New(l)

	doc := p.Parse()
	if doc == nil {
		t.Fatal("fail to parse DOM")
	}
	if doc.FirstChild == nil {
		t.Fatal("first child missing")
	}
	n := doc.FirstChild
	for i, tt := range tests {
		if n.Tag != tt.Tag {
			t.Fatalf("tests[%d]: invalid tag. expected=%q got=%q", i, tt.Tag, n.Tag)
		}
	}
}
