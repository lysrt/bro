package css

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		.name, a {background-color: #FF00FF; padding: 12em;}
		#figure {
			color : rgb(128, 0, 0);
			margin: 0;
		}
		/*    comment example*/
		>
	`

	tests := []CSSToken{
		{DOT, "."},
		{IDENTIFIER, "name"},
		{COMMA, ","},
		{IDENTIFIER, "a"},
		{LBRACE, "{"},
		{IDENTIFIER, "background-color"},
		{COLON, ":"},
		{HASH, "#"},
		{IDENTIFIER, "FF00FF"},
		{SEMICOLON, ";"},
		{IDENTIFIER, "padding"},
		{COLON, ":"},
		{NUMBER, "12"},
		{IDENTIFIER, "em"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{HASH, "#"},
		{IDENTIFIER, "figure"},
		{LBRACE, "{"},
		{IDENTIFIER, "color"},
		{COLON, ":"},
		{IDENTIFIER, "rgb"},
		{LPARENTHESIS, "("},
		{NUMBER, "128"},
		{COMMA, ","},
		{NUMBER, "0"},
		{COMMA, ","},
		{NUMBER, "0"},
		{RPARENTHESIS, ")"},
		{SEMICOLON, ";"},
		{IDENTIFIER, "margin"},
		{COLON, ":"},
		{NUMBER, "0"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{COMMENT, "comment example"},
		{ILLEGAL, ">"},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q", i, tt.Type, tok.Type)
		}

		if tok.Litteral != tt.Litteral {
			t.Fatalf("tests[%d] - wrong token litteral. expected=%q, got=%q", i, tt.Litteral, tok.Litteral)
		}
	}
}
