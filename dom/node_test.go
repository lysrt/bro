package dom

import (
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

// TODO: Simple test, just checking the number of children, maybe a bit light...
func TestNodeChildren(t *testing.T) {
	tests := []struct {
		input    *html.Node
		expected int
	}{
		{htmlParseSnippet(t, "<p></p>"), 0},
		{htmlParseSnippet(t, "<p><span/></p>"), 1},
		{htmlParseSnippet(t, "<div><p></p><a></a><p><b></b></p></div>"), 3},
		{htmlParseSnippet(t, "<p><div></div></p>"), 0}, // Special case with div in p
	}

	for _, tt := range tests {
		actual := NodeChildren(tt.input)
		if len(actual) != tt.expected {
			t.Errorf("NodeChildren: actual %v, expected %v", actual, tt.expected)
		}
	}
}
