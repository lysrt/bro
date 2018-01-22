package parser

import (
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

func TestDOMElement(t *testing.T) {
	input := `<a></a><b></b><c></c>`
	tests := []*dom.Element{
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
	if len(doc.Children) != len(tests) {
		t.Fatalf("invalid children count. expected=%d got=%d", len(tests), len(doc.Children))
	}
	for i, tt := range tests {
		e, ok := doc.Children[i].(*dom.Element)
		if !ok {
			t.Fatalf("tests[%d]: invalid node type. expected=%T got=%T", i, tt, doc.Children[i])
		}
		if e.Tag != tt.Tag {
			t.Fatalf("tests[%d]: invalid tag. expected=%q got=%q", i, tt.Tag, e.Tag)
		}
	}
}
