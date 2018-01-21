package lexer

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	Identifier = "Identifier" // id, class, charset, ...

	Bang          = "!"
	Equal         = "="
	Slash         = "/"
	String        = "String"
	Text          = "Text"
	FAngleBracket = ">" // forward angle bracket
	BAngleBracket = "<" // backward angle bracket
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
