package layout

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/style"

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

	blockStyle := css.NewParser(strings.NewReader(blockCSS)).ParseStylesheet()
	inlineStyle := css.NewParser(strings.NewReader(blockCSS)).ParseStylesheet()
	mixedStyle := css.NewParser(strings.NewReader(blockCSS)).ParseStylesheet()

	type args struct {
		styleTree *style.StyledNode
	}
	tests := []struct {
		name string
		args args
		want *LayoutBox
	}{
		{
			name: "Block Layout",
			args: args{style.GenerateStyleTree(node, blockStyle)},
			want: &LayoutBox{
				BoxType: BlockNode,
				Children: []*LayoutBox{
					{BoxType: BlockNode},
					{BoxType: BlockNode},
					{BoxType: BlockNode},
					{BoxType: BlockNode},
				},
			},
		},
		{
			name: "Inline Layout",
			args: args{style.GenerateStyleTree(node, inlineStyle)},
			want: &LayoutBox{
				BoxType: BlockNode,
				Children: []*LayoutBox{
					{BoxType: InlineNode},
					{BoxType: InlineNode},
					{BoxType: InlineNode},
					{BoxType: InlineNode},
				},
			},
		},
		{
			name: "Mixed Layout",
			args: args{style.GenerateStyleTree(node, mixedStyle)},
			want: &LayoutBox{
				BoxType: BlockNode,
				Children: []*LayoutBox{
					{BoxType: BlockNode},
					{BoxType: InlineNode},
					{BoxType: InlineNode},
					{BoxType: BlockNode},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			got := GenerateLayoutTree(tt.args.styleTree)
			if got.BoxType != want.BoxType {
				node := got.StyledNode.Node.Data
				t.Fatalf("node %q expects %q got %q", node, want.BoxType, got.BoxType)
			}
			for i := range got.Children {
				got := got.Children[i]
				want := want.Children[i]
				if got.BoxType != want.BoxType {
					node := got.StyledNode.Node.Data
					t.Fatalf("node %q expects %q got %q", node, want.BoxType, got.BoxType)
				}
			}
		})
	}
}
