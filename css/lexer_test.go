package css

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := ".name, a { background-color: #FFFF00FF; padding: 12em;}"

	tests := []struct {
		expectedType     CSSTokenType
		expectedLitteral string
	}{
		{DOT, "."},
		{IDENTIFIER, "name"},
		{COMMA, ","},
		{IDENTIFIER, "a"},
		{LBRACE, "{"},
		{IDENTIFIER, "background-color"},
		{COLON, ":"},
		{HASH, "#"},
		{IDENTIFIER, "FFFF00FF"},
		{SEMICOLON, ";"},
		{IDENTIFIER, "padding"},
		{COLON, ":"},
		{NUMBER, "12"},
		{IDENTIFIER, "em"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Litteral != tt.expectedLitteral {
			t.Fatalf("tests[%d] - wrong token litteral. expected=%q, got=%q", i, tt.expectedLitteral, tok.Litteral)
		}
	}
}
