package main

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/dom"

	"github.com/lysrt/bro/css"
	"golang.org/x/net/html"
)

var blockHTML = `<container>
	<a></a>
	<b></b>
	<c></c>
	<d></d>
</container>`

var (
	blockCSS  = `a, b, c, d { display: block; }`
	inlineCSS = `a, b, c, d { display: inline; }`
	mixedCSS  = `a, d { display: block; } b, c { display: inline; }`
)

func TestGenerateLayoutTree(t *testing.T) {
	node, err := html.Parse(strings.NewReader(blockHTML))
	if err != nil {
		t.Fatal("fail to parse DOM:", err)
	}
	// go through ??? -> html -> body -> container
	node = dom.NodeLastElementChild(node)
	node = dom.NodeLastElementChild(node)
	node = dom.NodeLastElementChild(node)
	t.Log(node.Data)

	blockStyle := css.ParseCSS(strings.NewReader(blockCSS))
	inlineStyle := css.ParseCSS(strings.NewReader(blockCSS))
	mixedStyle := css.ParseCSS(strings.NewReader(blockCSS))

	type args struct {
		styleTree *StyledNode
	}
	tests := []struct {
		name string
		args args
		want *LayoutBox
	}{
		{
			name: "Block Layout",
			args: args{GenerateStyleTree(node, blockStyle)},
			want: &LayoutBox{
				boxType: BlockNode,
				children: []*LayoutBox{
					{boxType: BlockNode},
					{boxType: BlockNode},
					{boxType: BlockNode},
					{boxType: BlockNode},
				},
			},
		},
		{
			name: "Inline Layout",
			args: args{GenerateStyleTree(node, inlineStyle)},
			want: &LayoutBox{
				boxType: BlockNode,
				children: []*LayoutBox{
					{boxType: InlineNode},
					{boxType: InlineNode},
					{boxType: InlineNode},
					{boxType: InlineNode},
				},
			},
		},
		{
			name: "Mixed Layout",
			args: args{GenerateStyleTree(node, mixedStyle)},
			want: &LayoutBox{
				boxType: BlockNode,
				children: []*LayoutBox{
					{boxType: BlockNode},
					{boxType: InlineNode},
					{boxType: InlineNode},
					{boxType: BlockNode},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			got := GenerateLayoutTree(tt.args.styleTree)
			if got.boxType != want.boxType {
				node := got.styledNode.Node.Data
				t.Fatalf("node %q expects %q got %q", node, want.boxType, got.boxType)
			}
			for i := range got.children {
				got := got.children[i]
				want := want.children[i]
				if got.boxType != want.boxType {
					node := got.styledNode.Node.Data
					t.Fatalf("node %q expects %q got %q", node, want.boxType, got.boxType)
				}
			}
		})
	}
}
