package lexer

import "fmt"

const (
	TokenError = "Error"
	TokenEOF   = "EOF"

	TokenIdent   = "Identifier" // id, class, charset, ...
	TokenString  = "String"
	TokenText    = "Text"
	TokenComment = "Comment"

	TokenBang     = "!"
	TokenEqual    = "="
	TokenSlash    = "/"
	TokenRBracket = ">" // forward angle bracket
	TokenLBracket = "<" // backward angle bracket
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string

	Position     int
	Line         int
	LinePosition int
}

func (t Token) String() string {
	switch t.Type {
	case TokenError:
		return fmt.Sprintf("error at line %d (%d): %s", t.Line, t.LinePosition, t.Literal)
	case TokenIdent, TokenString, TokenText, TokenComment:
		if len(t.Literal) > 10 {
			return fmt.Sprintf("%.10q...", t.Literal)
		}
		return fmt.Sprintf("%q", t.Literal)
	default:
		return string(t.Type)
	}
}
