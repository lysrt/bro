package main

import (
	"reflect"
	"testing"
)

var (
	blockBoxes  *StyledNode
	inlineBoxes *StyledNode
	mixedBoxes  *StyledNode
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

func init() {

}

func TestGenerateLayoutTree(t *testing.T) {
	type args struct {
		styleTree *StyledNode
	}
	tests := []struct {
		name string
		args args
		want *LayoutBox
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateLayoutTree(tt.args.styleTree); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateLayoutTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
