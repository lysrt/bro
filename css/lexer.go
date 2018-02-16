package css

import (
	"strings"
)

type CSSTokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	COMMENT = "COMMENT"

	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"

	STAR         = "*"
	DOT          = "."
	COMMA        = ","
	COLON        = ":"
	SEMICOLON    = ";"
	HASH         = "#"
	LBRACE       = "{"
	RBRACE       = "}"
	LPARENTHESIS = "("
	RPARENTHESIS = ")"
)

type CSSToken struct {
	Type     CSSTokenType
	Litteral string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.char = l.peekChar()
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() CSSToken {
	var tok CSSToken

	l.skipWhitespace()
	switch l.char {
	case '*':
		tok = newToken(STAR, l.char)
	case '.':
		tok = newToken(DOT, l.char)
	case ',':
		tok = newToken(COMMA, l.char)
	case ':':
		tok = newToken(COLON, l.char)
	case ';':
		tok = newToken(SEMICOLON, l.char)
	case '#':
		tok = newToken(HASH, l.char)
	case '{':
		tok = newToken(LBRACE, l.char)
	case '}':
		tok = newToken(RBRACE, l.char)
	case '(':
		tok = newToken(LPARENTHESIS, l.char)
	case ')':
		tok = newToken(RPARENTHESIS, l.char)
	case '/':
		next := l.peekChar()
		if next != '*' {
			tok = newToken(ILLEGAL, l.char)
		} else {
			comment := l.readComment()
			tok.Litteral = strings.TrimSpace(comment)
			tok.Type = COMMENT
			return tok
		}
	case 0:
		tok.Litteral = ""
		tok.Type = EOF
	default:
		if isLetter(l.char) {
			tok.Litteral = l.readIdentifier()
			tok.Type = IDENTIFIER
			return tok
		} else if isDigit(l.char) {
			tok.Type = NUMBER
			tok.Litteral = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.char)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType CSSTokenType, char byte) CSSToken {
	return CSSToken{Type: tokenType, Litteral: string(char)}
}

func (l *Lexer) readComment() string {
	l.readChar()
	l.readChar()
	position := l.position
	for l.char != '*' || l.input[l.readPosition] != '/' {
		l.readChar()
	}
	comment := l.input[position:l.position]

	l.readChar()
	l.readChar()

	return comment
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifierPart(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isNumberPart(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char == '-'
}

func isIdentifierPart(char byte) bool {
	return isLetter(char) || '0' <= char && char <= '9'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isNumberPart(char byte) bool {
	return isDigit(char) || char == '.'
}
