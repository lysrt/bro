package parser

import (
	"strings"
	"testing"

	"github.com/lysrt/bro/dom"
	"github.com/lysrt/bro/dom/lexer"
)

const sample = ""

func BenchmarkGoHTMLParser(b *testing.B) {
	r := strings.NewReader(sample)

	for n := 0; n < b.N; n++ {
		_, _ = dom.ParseHTML(r)
	}
}

func BenchmarkCustomHTMLParser(b *testing.B) {
	l := lexer.New(sample)
	p := New(l)

	for n := 0; n < b.N; n++ {
		_ = p.Parse()
	}
}
