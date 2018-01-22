package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `<!doctype html>
<html class="app--red" data-awesome='true'>
	Hello world!
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
		{Type: TokenText, Literal: "\n"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenIdent, Literal: "html"},
		{Type: TokenIdent, Literal: "class"},
		{Type: TokenEqual, Literal: "="},
		{Type: TokenString, Literal: "app--red"},
		{Type: TokenIdent, Literal: "data-awesome"},
		{Type: TokenEqual, Literal: "="},
		{Type: TokenString, Literal: "true"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenText, Literal: "\n\tHello world!\n"},
		{Type: TokenRBracket, Literal: "<"},
		{Type: TokenSlash, Literal: "/"},
		{Type: TokenIdent, Literal: "html"},
		{Type: TokenLBracket, Literal: ">"},
		{Type: TokenEOF, Literal: ""},
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

func TestNextToken_position(t *testing.T) {
	input := `<apple class="apple--red">
</apple>`

	tests := []Token{
		{Literal: "<", Position: 0, Line: 0, LinePosition: 0},
		{Literal: "apple", Position: 1, Line: 0, LinePosition: 1},
		{Literal: "class", Position: 7, Line: 0, LinePosition: 7},
		{Literal: "=", Position: 12, Line: 0, LinePosition: 12},
		{Literal: "apple--red", Position: 13, Line: 0, LinePosition: 13},
		{Literal: ">", Position: 25, Line: 0, LinePosition: 25},
		{Literal: "\n", Position: 26, Line: 0, LinePosition: 26},
		{Literal: "<", Position: 27, Line: 1, LinePosition: 0},
		{Literal: "/", Position: 28, Line: 1, LinePosition: 1},
		{Literal: "apple", Position: 29, Line: 1, LinePosition: 2},
		{Literal: ">", Position: 34, Line: 1, LinePosition: 7},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Position != tt.Position {
			t.Fatalf("tests[%d]: bad position. expected=%d got=%d", i, tt.Position, tok.Position)
		}
		if tok.Line != tt.Line {
			t.Fatalf("tests[%d]: bad line. expected=%d got=%d", i, tt.Line, tok.Line)
		}
		if tok.LinePosition != tt.LinePosition {
			t.Fatalf("tests[%d]: bad line position. expected=%d got=%d", i, tt.LinePosition, tok.LinePosition)
		}
	}
}
