package layout

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/html"
	"github.com/lysrt/bro/html/lexer"
	"github.com/lysrt/bro/html/parser"
	"github.com/lysrt/bro/style"

	"github.com/lysrt/bro/css"
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
	l := lexer.New(blockHTML)
	p := parser.New(l)
	node := p.Parse()
	if errors := p.Errors(); len(errors) > 0 {
		t.Fatal(errors)
	}
	// go through ??? -> html -> body -> container
	node = html.NodeLastElementChild(node)
	node = html.NodeLastElementChild(node)
	node = html.NodeLastElementChild(node)
	t.Log(node.Tag)

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
				node := got.StyledNode.Node.Tag
				t.Fatalf("node %q expects %q got %q", node, want.BoxType, got.BoxType)
			}
			for i := range got.Children {
				got := got.Children[i]
				want := want.Children[i]
				if got.BoxType != want.BoxType {
					node := got.StyledNode.Node.Tag
					t.Fatalf("node %q expects %q got %q", node, want.BoxType, got.BoxType)
				}
			}
		})
	}
}
