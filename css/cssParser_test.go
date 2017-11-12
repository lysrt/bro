package css

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/parser"
)

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
	{"#AAEE120C", Color{A: 170, R: 238, G: 18, B: 12}},
}

func Test_parseColor(t *testing.T) {
	for _, tt := range ColorTests {
		var s parser.Scanner
		s.Init(strings.NewReader(tt.input))

		actual := parseColor(&s)

		if actual != tt.expected {
			t.Errorf("%s - expected: %v actual:  %v", tt.input, tt.expected, actual)
		}
	}
}
