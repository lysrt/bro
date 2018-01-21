package lexer

const (
	TokenIllegal = "Illegal"
	TokenEOF     = "EOF"

	TokenIdent  = "Identifier" // id, class, charset, ...
	TokenString = "String"
	TokenText   = "Text"

	TokenBang     = "!"
	TokenEqual    = "="
	TokenSlash    = "/"
	TokenLBracket = ">" // forward angle bracket
	TokenRBracket = "<" // backward angle bracket
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
