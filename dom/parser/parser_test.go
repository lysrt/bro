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
	}
}
