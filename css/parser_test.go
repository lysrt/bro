package css

import (
	"strings"
	"testing"
)

var SelectorTests = []struct {
	input    string
	expected Selector
	isErr    bool
}{
	{"#id", Selector{ID: "id"}, false},
	{".class", Selector{Classes: []string{"class"}}, false},
	{"tag", Selector{TagName: "tag"}, false},
	{"..", Selector{}, true},
	{"#/", Selector{}, true},
}

func TestSelector(t *testing.T) {
	for _, tt := range SelectorTests {
		p := NewParser(strings.NewReader(tt.input))

		selector := p.parseSelector()

		actual := selector
		if actual.ID != tt.expected.ID {
			t.Fatalf("%s - expected: %v actual:  %v", tt.input, tt.expected, actual)
		}
		if len(actual.Classes) != len(tt.expected.Classes) {
			t.Fatalf("%s - expected: %v actual:  %v", tt.input, tt.expected, actual)
		}
		if actual.TagName != tt.expected.TagName {
			t.Fatalf("%s - expected: %v actual:  %v", tt.input, tt.expected, actual)
		}
	}
}
func TestSelectors(t *testing.T) {
	p := NewParser(strings.NewReader("#id, .class, tag"))

	selectors := p.parseSelectors()
	if len(selectors) != 3 {
		t.Fatal("wrong number of selectors")
	}
}

func TestKeywordDeclaration(t *testing.T) {
	p := NewParser(strings.NewReader("margin: auto;"))

	declaration := p.parseDeclaration()
	if declaration.Name != "margin" {
		t.Fatal("wrong keyword declaration name")
	}

	expected := Value{Keyword: "auto"}
	if declaration.Value != expected {
		t.Fatal("wrong keyword declaration value")
	}
}

func TestColorDeclaration(t *testing.T) {
	p := NewParser(strings.NewReader("color: #FFFFFF;"))

	declaration := p.parseDeclaration()
	if declaration.Name != "color" {
		t.Fatal("wrong color declaration name")
	}

	expected := Value{Color: Color{"", 255, 255, 255, 255}}
	if declaration.Value != expected {
		t.Fatal("wrong color declaration value")
	}
}

func TestLengthDeclaration(t *testing.T) {
	p := NewParser(strings.NewReader("size: 50.5 px;"))

	declaration := p.parseDeclaration()
	if declaration.Name != "size" {
		t.Fatal("wrong length declaration name")
	}

	expected := Value{Length: Length{Quantity: 50.5, Unit: Px}}
	if declaration.Value != expected {
		t.Fatal("wrong length declaration value")
	}
}

var ColorTests = []struct {
	input    string
	expected Color
}{
	{"#FFFFFF", Color{A: 255, R: 255, G: 255, B: 255}},
	{"#ffffff", Color{A: 255, R: 255, G: 255, B: 255}},
	{"#fff", Color{A: 255, R: 255, G: 255, B: 255}},
	{"#000000", Color{A: 255, R: 0, G: 0, B: 0}},
	{"#DD0001", Color{A: 255, R: 221, G: 0, B: 1}},
	{"#abc", Color{A: 255, R: 170, G: 187, B: 204}},
}

func TestColor(t *testing.T) {
	for _, tt := range ColorTests {
		p := NewParser(strings.NewReader(tt.input))

		actual := p.parseColor()

		if actual != tt.expected {
			t.Errorf("%s - expected: %v actual:  %v", tt.input, tt.expected, actual)
		}
	}
}
