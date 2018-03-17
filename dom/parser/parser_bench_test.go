package parser

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/dom/lexer"
	"golang.org/x/net/html"
)

const sample = `<div class="a"><div class="b"><div class="c"><div class="d"><div class="e"><div class="f"><div class="g"><div class="h"></div></div></div></div></div></div></div></div>`

func BenchmarkGoHTMLParser(b *testing.B) {
	r := strings.NewReader(sample)

	for n := 0; n < b.N; n++ {
		_, err := html.Parse(r)
		if err != nil {
			b.Error("Parsing error", err)
		}
	}
}

func BenchmarkCustomHTMLParser(b *testing.B) {
	l := lexer.New(sample)
	p := New(l)

	for n := 0; n < b.N; n++ {
		_ = p.Parse()
		if len(p.Errors()) > 0 {
			b.Error("Parsing errors:")
		}
	}
}
