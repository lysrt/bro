package main

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

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
				n:        &html.Node{Data: "p"},
				selector: Selector{TagName: "p"},
			},
			want: true,
		},
		{
			name: "dont match tag",
			args: args{
				n:        &html.Node{Data: "div"},
				selector: Selector{TagName: "p"},
			},
			want: false,
		},
		{
			name: "match id",
			args: args{
				n: &html.Node{
					Data: "div",
					Attr: []html.Attribute{{Key: "id", Val: "bloup"}},
				},
				selector: Selector{ID: "bloup"},
			},
			want: true,
		},
		{
			name: "match class",
			args: args{
				n: &html.Node{
					Data: "div",
					Attr: []html.Attribute{{Key: "class", Val: "bloup blip"}},
				},
				selector: Selector{Class: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex",
			args: args{
				n: &html.Node{
					Data: "div",
					Attr: []html.Attribute{{Key: "class", Val: "bloup blip"}},
				},
				selector: Selector{TagName: "div", Class: []string{"blip"}},
			},
			want: true,
		},
		{
			name: "match complex (invalid class)",
			args: args{
				n: &html.Node{
					Data: "div",
					Attr: []html.Attribute{{Key: "class", Val: "bloup blip"}},
				},
				selector: Selector{TagName: "div", Class: []string{"arf"}},
			},
			want: false,
		},
		{
			name: "match complex (invalid tag)",
			args: args{
				n: &html.Node{
					Data: "div",
					Attr: []html.Attribute{{Key: "class", Val: "bloup blip"}},
				},
				selector: Selector{TagName: "p", Class: []string{"bloup", "blip"}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchSelector(tt.args.n, tt.args.selector); got != tt.want {
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
