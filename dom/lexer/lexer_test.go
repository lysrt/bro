package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `<!doctype html>
<html class="app--red">
<body class='app__container'>
	Hello world!
</body>
</html>`

	tests := []struct {
		Type    TokenType
		Literal string
	}{
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenBang, Literal: "!"},
		{Type: TokenIdent, Literal: "doctype"},
		{Type: TokenIdent, Literal: "html"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenIdent, Literal: "html"},
		{Type: TokenIdent, Literal: "class"},
		{Type: TokenEqual, Literal: "="},
		{Type: TokenString, Literal: "app--red"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenIdent, Literal: "body"},
		{Type: TokenIdent, Literal: "class"},
		{Type: TokenEqual, Literal: "="},
		{Type: TokenString, Literal: "app__container"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenText, Literal: "\n\tHello world!\n"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenSlash, Literal: "/"},
		{Type: TokenIdent, Literal: "body"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenSlash, Literal: "/"},
		{Type: TokenIdent, Literal: "html"},
		{Type: TokenLBracket, Literal: ">"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.Type {
			t.Fatalf("tests[%d]: bad token type. expected=%q got=%q", i, tt.Type, tok.Type)
		}
		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d]: bad literal. expected=%q got=%q", i, tt.Literal, tok.Literal)
		}
	}
}
