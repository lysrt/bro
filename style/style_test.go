package style

import (
	"reflect"
	"testing"

	"github.com/lysrt/bro/css"
	"github.com/lysrt/bro/html"
	"github.com/lysrt/bro/html/lexer"
	"github.com/lysrt/bro/html/parser"
)

func htmlParseSnippet(t *testing.T, data string) *html.Node {
	l := lexer.New(data)
	p := parser.New(l)
	n := p.Parse()
	if errors := p.Errors(); len(errors) > 0 {
		t.Fatal(errors)
	}
	// return the first element of the body.
	return n.LastChild.FirstChild
}

func Test_matchSelector(t *testing.T) {
	type args struct {
		n        *html.Node
		selector css.Selector
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match tag",
			args: args{
				n:        htmlParseSnippet(t, "<p></p>"),
				selector: css.Selector{TagName: "p"},
			},
			want: true,
		},
		{
			name: "dont match tag",
			args: args{
				n:        htmlParseSnippet(t, "<div></div>"),
				selector: css.Selector{TagName: "p"},
			},
			want: false,
		},
		{
			name: "match id",
			args: args{
				n:        htmlParseSnippet(t, `<div id="bloup"></div>`),
				selector: css.Selector{ID: "bloup"},
			},
			want: true,
		},
		{
			name: "match class",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: css.Selector{Classes: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: css.Selector{TagName: "div", Classes: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex (invalid class)",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: css.Selector{TagName: "div", Classes: []string{"arf"}},
			},
			want: false,
		},
		{
			name: "match complex (invalid tag)",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: css.Selector{TagName: "p", Classes: []string{"bloup", "blip"}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchSelector(tt.args.n, tt.args.selector); got != tt.want {
				// TODO: write HTML into buffer
				//buf := bytes.Buffer{}
				//html.Render(&buf, tt.args.n)
				//t.Log(buf.String())
				t.Errorf("matchSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchRule(t *testing.T) {
	type args struct {
		n *html.Node
		r css.Rule
	}
	tests := []struct {
		name     string
		args     args
		wantRule css.Rule
		wantOk   bool
	}{
		{
			name: "valid",
			args: args{
				n: &html.Node{Tag: "p"},
				r: css.Rule{Selectors: []css.Selector{{TagName: "a"}, {TagName: "p"}}},
			},
			wantRule: css.Rule{Selectors: []css.Selector{{TagName: "a"}, {TagName: "p"}}},
			wantOk:   true,
		},
		{
			name: "invalid",
			args: args{
				n: &html.Node{Tag: "p"},
				r: css.Rule{Selectors: []css.Selector{{TagName: "a"}, {ID: "bloup"}}},
			},
			wantRule: css.Rule{},
			wantOk:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, gotOk := matchRule(tt.args.n, tt.args.r)
			if !reflect.DeepEqual(gotM.Rule, tt.wantRule) {
				t.Errorf("matchRule() gotM = %v, want %v", gotM.Rule, tt.wantRule)
			}
			if gotOk != tt.wantOk {
				t.Errorf("matchRule() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_specifiedValues(t *testing.T) {
	node := htmlParseSnippet(t, `<ul id="list" class="bullet center">
		<li class="first all"></li>
		<li class="all"></li>
		<li id="green-point" class="last all"></li>
	</ul>`)
	style := &css.Stylesheet{
		Rules: []css.Rule{
			{
				Selectors: []css.Selector{
					{TagName: "ul"},
				},
				Declarations: []css.Declaration{
					{Name: "background-color", Value: css.Value{Color: css.Color{Name: "red"}}},
					{Name: "width", Value: css.Value{Length: css.Length{Quantity: 900, Unit: css.Px}}},
				},
			},
			{
				Selectors: []css.Selector{
					{ID: "list"},
				},
				Declarations: []css.Declaration{
					{Name: "background-color", Value: css.Value{Color: css.Color{Name: "blue"}}},
				},
			},
			{
				Selectors: []css.Selector{
					{Classes: []string{"all"}},
				},
				Declarations: []css.Declaration{
					{Name: "font-size", Value: css.Value{Length: css.Length{Quantity: 12, Unit: css.Px}}},
				},
			},
			{
				Selectors: []css.Selector{
					{Classes: []string{"first"}},
				},
				Declarations: []css.Declaration{
					{Name: "color", Value: css.Value{Color: css.Color{Name: "blue"}}},
				},
			},
			{
				Selectors: []css.Selector{
					{Classes: []string{"last"}},
				},
				Declarations: []css.Declaration{
					{Name: "color", Value: css.Value{Color: css.Color{Name: "red"}}},
				},
			},
			{
				Selectors: []css.Selector{
					{ID: "green-point"},
				},
				Declarations: []css.Declaration{
					{Name: "color", Value: css.Value{Color: css.Color{Name: "green"}}},
				},
			},
		},
	}
	type args struct {
		n          *html.Node
		stylesheet *css.Stylesheet
	}
	tests := []struct {
		name string
		args args
		want PropertyMap
	}{
		{
			name: "ul#list",
			args: args{
				n:          node,
				stylesheet: style,
			},
			want: PropertyMap{
				"background-color": css.Value{Color: css.Color{Name: "blue"}},
				"width":            css.Value{Length: css.Length{Quantity: 900, Unit: css.Px}},
			},
		},
		{
			name: ".first.all",
			args: args{
				n:          html.NodeFirstElementChild(node),
				stylesheet: style,
			},
			want: PropertyMap{
				"color":     css.Value{Color: css.Color{Name: "blue"}},
				"font-size": css.Value{Length: css.Length{Quantity: 12, Unit: css.Px}},
			},
		},
		{
			name: "#green-point.last.all",
			args: args{
				n:          html.NodeLastElementChild(node),
				stylesheet: style,
			},
			want: PropertyMap{
				"color":     css.Value{Color: css.Color{Name: "green"}},
				"font-size": css.Value{Length: css.Length{Quantity: 12, Unit: css.Px}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := specifiedValues(tt.args.n, tt.args.stylesheet); !reflect.DeepEqual(got, tt.want) {
				// TODO: write HTML into buffer
				//buf := bytes.Buffer{}
				//html.Render(&buf, tt.args.n)
				//t.Log(buf.String())
				t.Errorf("specifiedValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
