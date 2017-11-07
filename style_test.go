package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func htmlParseSnippet(t *testing.T, data string) *html.Node {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal("fail to parse HTML:", err)
	}
	// return the first element of the body.
	return n.FirstChild.LastChild.FirstChild
}

func Test_matchSelector(t *testing.T) {
	type args struct {
		n        *html.Node
		selector Selector
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
				selector: Selector{TagName: "p"},
			},
			want: true,
		},
		{
			name: "dont match tag",
			args: args{
				n:        htmlParseSnippet(t, "<div></div>"),
				selector: Selector{TagName: "p"},
			},
			want: false,
		},
		{
			name: "match id",
			args: args{
				n:        htmlParseSnippet(t, `<div id="bloup"></div>`),
				selector: Selector{ID: "bloup"},
			},
			want: true,
		},
		{
			name: "match class",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: Selector{Classes: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: Selector{TagName: "div", Classes: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex (invalid class)",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: Selector{TagName: "div", Classes: []string{"arf"}},
			},
			want: false,
		},
		{
			name: "match complex (invalid tag)",
			args: args{
				n:        htmlParseSnippet(t, `<div class="bloup blip"></div>`),
				selector: Selector{TagName: "p", Classes: []string{"bloup", "blip"}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchSelector(tt.args.n, tt.args.selector); got != tt.want {
				buf := bytes.Buffer{}
				html.Render(&buf, tt.args.n)
				t.Log(buf.String())
				t.Errorf("matchSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchRule(t *testing.T) {
	type args struct {
		n *html.Node
		r Rule
	}
	tests := []struct {
		name     string
		args     args
		wantRule Rule
		wantOk   bool
	}{
		{
			name: "valid",
			args: args{
				n: &html.Node{Data: "p"},
				r: Rule{Selectors: []Selector{{TagName: "a"}, {TagName: "p"}}},
			},
			wantRule: Rule{Selectors: []Selector{{TagName: "a"}, {TagName: "p"}}},
			wantOk:   true,
		},
		{
			name: "invalid",
			args: args{
				n: &html.Node{Data: "p"},
				r: Rule{Selectors: []Selector{{TagName: "a"}, {ID: "bloup"}}},
			},
			wantRule: Rule{},
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
	style := Stylesheet{
		Rules: []Rule{
			{
				Selectors: []Selector{
					{TagName: "ul"},
				},
				Declarations: []Declaration{
					{Name: "background-color", Value: Value{Color: Color{Name: "red"}}},
					{Name: "width", Value: Value{Length: Length{Quantity: 900, Unit: Px}}},
				},
			},
			{
				Selectors: []Selector{
					{ID: "list"},
				},
				Declarations: []Declaration{
					{Name: "background-color", Value: Value{Color: Color{Name: "blue"}}},
				},
			},
			{
				Selectors: []Selector{
					{Classes: []string{"all"}},
				},
				Declarations: []Declaration{
					{Name: "font-size", Value: Value{Length: Length{Quantity: 12, Unit: Px}}},
				},
			},
			{
				Selectors: []Selector{
					{Classes: []string{"first"}},
				},
				Declarations: []Declaration{
					{Name: "color", Value: Value{Color: Color{Name: "blue"}}},
				},
			},
			{
				Selectors: []Selector{
					{Classes: []string{"last"}},
				},
				Declarations: []Declaration{
					{Name: "color", Value: Value{Color: Color{Name: "red"}}},
				},
			},
			{
				Selectors: []Selector{
					{ID: "green-point"},
				},
				Declarations: []Declaration{
					{Name: "color", Value: Value{Color: Color{Name: "green"}}},
				},
			},
		},
	}
	type args struct {
		n          *html.Node
		stylesheet Stylesheet
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
				"background-color": Value{Color: Color{Name: "blue"}},
				"width":            Value{Length: Length{Quantity: 900, Unit: Px}},
			},
		},
		{
			name: ".first.all",
			args: args{
				n:          NodeFirstElementChild(node),
				stylesheet: style,
			},
			want: PropertyMap{
				"color":     Value{Color: Color{Name: "blue"}},
				"font-size": Value{Length: Length{Quantity: 12, Unit: Px}},
			},
		},
		{
			name: "#green-point.last.all",
			args: args{
				n:          NodeLastElementChild(node),
				stylesheet: style,
			},
			want: PropertyMap{
				"color":     Value{Color: Color{Name: "green"}},
				"font-size": Value{Length: Length{Quantity: 12, Unit: Px}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := specifiedValues(tt.args.n, tt.args.stylesheet); !reflect.DeepEqual(got, tt.want) {
				buf := bytes.Buffer{}
				html.Render(&buf, tt.args.n)
				t.Log(buf.String())
				t.Errorf("specifiedValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
