package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `<!doctype html>
<html class="app--red">
<body class='app__container'>
</body>
</html>`

	tests := []struct {
		Type    TokenType
		Literal string
	}{
		{Type: BAngleBracket, Literal: "<"},
		{Type: Bang, Literal: "!"},
		{Type: Identifier, Literal: "doctype"},
		{Type: Identifier, Literal: "html"},
		{Type: FAngleBracket, Literal: ">"},
		{Type: BAngleBracket, Literal: "<"},
		{Type: Identifier, Literal: "html"},
		{Type: Identifier, Literal: "class"},
		{Type: Equal, Literal: "="},
		{Type: String, Literal: "app--red"},
		{Type: FAngleBracket, Literal: ">"},
		{Type: BAngleBracket, Literal: "<"},
		{Type: Identifier, Literal: "body"},
		{Type: Identifier, Literal: "class"},
		{Type: Equal, Literal: "="},
		{Type: String, Literal: "app__container"},
		{Type: FAngleBracket, Literal: ">"},
		{Type: BAngleBracket, Literal: "<"},
		{Type: Slash, Literal: "/"},
		{Type: Identifier, Literal: "body"},
		{Type: FAngleBracket, Literal: ">"},
		{Type: BAngleBracket, Literal: "<"},
		{Type: Slash, Literal: "/"},
		{Type: Identifier, Literal: "html"},
		{Type: FAngleBracket, Literal: ">"},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.Type {
			t.Fatalf("bad token type. expected=%q got=%q", tt.Type, tok.Type)
		}
		if tok.Literal != tt.Literal {
			t.Fatalf("bad literal. expected=%q got=%q", tt.Literal, tok.Literal)
		}
	}
}
